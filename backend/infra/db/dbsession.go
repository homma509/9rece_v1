package db

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
)

// Session DB接続の構造体
type Session struct {
	db        *dynamo.DB
	table     *dynamo.Table
	config    *aws.Config
	tableName string
}

// Resource DBリソースの構造体
type Resource interface {
	EntityName() string
	GetPK() string
	SetPK()
	GetSK() string
	SetSK()
	GetOwnerID() string
	SetOwnerID(string)
	GetVersion() uint64
	SetVersion(uint64)
	GetCreatedAt() time.Time
	SetCreatedAt(time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time)
}

// NewSession DB接続を生成します
func NewSession(config *aws.Config, tableName string) *Session {
	sess := &Session{
		config:    config,
		tableName: tableName,
	}
	return sess
}

func (s *Session) connect() error {
	if s.db == nil {
		sess, err := session.NewSession(s.config)
		if err != nil {
			return errors.WithStack(err)
		}
		s.db = dynamo.New(sess)
	}
	return nil
}

func (s *Session) connectTable() error {
	if s.table == nil {
		err := s.connect()
		if err != nil {
			return errors.WithStack(err)
		}

		table := s.db.Table(s.tableName)
		s.table = &table
	}
	return nil
}

// PutResource リソースをDBに登録します
func (s *Session) PutResource(r Resource) error {
	if s.isNewEntity(r) {
		return s.insertResource(r)
	}
	return s.updateResource(r)
}

// PutResources リソースのスライスをDBに登録します
func (s *Session) PutResources(rs []Resource) error {
	for _, r := range rs {
		if err := s.PutResource(r); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (s *Session) isNewEntity(r Resource) bool {
	return r.GetVersion() == 0
}

// GetResource リソースをDBから取得します
func (s *Session) GetResource(r Resource, ret interface{}) error {
	err := s.connectTable()
	if err != nil {
		return errors.WithStack(err)
	}

	err = s.table.
		Get("ID", r.GetPK()).
		Range("DataType", dynamo.Equal, r.GetSK()).
		One(ret)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Session) insertResource(r Resource) error {
	query, err := s.buildQueryInsert(r)
	if err != nil {
		return errors.WithStack(err)
	}

	err = query.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Session) updateResource(r Resource) error {
	query, err := s.buildQueryUpdate(r)
	if err != nil {
		return errors.WithStack(err)
	}

	err = query.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *Session) buildQueryInsert(r Resource) (*dynamo.Put, error) {
	err := s.connectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r.SetPK()
	r.SetSK()
	r.SetVersion(1)
	r.SetCreatedAt(time.Now())
	r.SetUpdatedAt(time.Now())

	query := s.table.Put(r)

	return query, nil
}

func (s *Session) buildQueryUpdate(r Resource) (*dynamo.Put, error) {
	err := s.connectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	oldVersion := r.GetVersion()

	r.SetVersion(oldVersion + 1)
	r.SetUpdatedAt(time.Now())

	query := s.table.Put(r)

	return query, nil
}
