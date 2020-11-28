package db

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/homma509/9rece/server/log"

	"github.com/pkg/errors"
)

const batchSize = 100

// Session DB接続の構造体
type Session struct {
	db        *dynamo.DB
	table     *dynamo.Table
	config    *aws.Config
	tableName string
}

// Resource DBリソースの構造体
type Resource interface {
	SetID()
	SetMetadata()
	SetCreatedAt(t time.Time)
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
	return s.insertResource(r)
}

// PutResources リソースのスライスをDBに登録します
func (s *Session) PutResources(rs []Resource) error {
	log.AppLogger.Info("Insert Count: ", len(rs))
	err := s.connectTable()
	if err != nil {
		return errors.WithStack(err)
	}

	is := make([]interface{}, len(rs))
	for i := range rs {
		rs[i].SetID()
		rs[i].SetMetadata()
		rs[i].SetCreatedAt(time.Now())
		is[i] = rs[i]
	}

	batches := make([][]interface{}, 0, (len(is)+batchSize-1)/batchSize)
	for batchSize < len(is) {
		is, batches = is[batchSize:], append(batches, is[0:batchSize:batchSize])
	}
	batches = append(batches, is)

	wrote := 0
	for _, batch := range batches {
		wrote, err = s.table.Batch().Write().Put(batch...).Run()
		if err != nil {
			return errors.WithStack(err)
		}
		wrote += wrote
	}
	if wrote != len(is) {
		errors.WithStack(fmt.Errorf("failed batch wrote %d ≠ %d(expected)", wrote, len(is)))
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

func (s *Session) buildQueryInsert(r Resource) (*dynamo.Put, error) {
	err := s.connectTable()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	r.SetID()
	r.SetMetadata()
	r.SetCreatedAt(time.Now())

	query := s.table.Put(r)

	return query, nil
}
