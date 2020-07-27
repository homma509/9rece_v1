package db

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/homma509/9rece/backend/domain/model"
	"github.com/pkg/errors"
)

// FacilityRepository 施設リポジトリの構造体
type FacilityRepository struct {
	sess *Session
}

// NewFacilityRepository 施設リポジトリを生成します
func NewFacilityRepository(sess *Session) *FacilityRepository {
	return &FacilityRepository{
		sess: sess,
	}
}

// Save 施設のスライスを登録します
func (r *FacilityRepository) Save(ctx context.Context, ps model.Facilities) error {
	rs := []Resource{}
	for _, p := range ps {
		rs = append(rs, r.newFacilityMapper(p))
	}

	err := r.sess.PutResources(rs)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *FacilityRepository) newFacilityMapper(p model.Facility) *FacilityMapper {
	new := &FacilityMapper{
		Facility: p,
	}

	var old FacilityMapper
	err := r.sess.GetResource(new, &old)
	if err == nil {
		new.SetPK()
		new.SetSK()
		new.SetOwnerID(old.OwnerID)
		new.SetCreatedAt(old.CreatedAt)
		new.SetVersion(old.Version)
	}

	return new
}

// FacilityMapper 施設モデルのリソースへのマッパー構造体
type FacilityMapper struct {
	model.Facility
	PK        string    `dynamo:"ID,hash"`
	SK        string    `dynamo:"DataType,range"`
	Version   uint64    `dynamo:"Version"`
	OwnerID   string    `dynamo:"OwnerID"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
	UpdatedAt time.Time `dynamo:"UpdatedAt"`
}

// EntityName Entity名の取得
func (m *FacilityMapper) EntityName() string {
	t := reflect.TypeOf(m.Facility)
	return t.Name()
}

// GetPK PKを取得します
func (m *FacilityMapper) GetPK() string {
	return m.FacilityID
}

// SetPK PKを設定します
func (m *FacilityMapper) SetPK() {
	m.PK = m.GetPK()
}

// GetSK SKを取得します
func (m *FacilityMapper) GetSK() string {
	return fmt.Sprintf("%sInfo", m.EntityName())
}

// SetSK SKを設定します
func (m *FacilityMapper) SetSK() {
	m.SK = m.GetSK()
}

// GetVersion バージョンを取得します
func (m *FacilityMapper) GetVersion() uint64 {
	return m.Version
}

// SetVersion Versionを設定します
func (m *FacilityMapper) SetVersion(v uint64) {
	m.Version = v
}

// GetOwnerID 所有者IDを取得します
func (m *FacilityMapper) GetOwnerID() string {
	return m.OwnerID
}

// SetOwnerID 所有者IDを設定します
func (m *FacilityMapper) SetOwnerID(id string) {
	m.OwnerID = id
}

// GetCreatedAt 登録日時を取得します
func (m *FacilityMapper) GetCreatedAt() time.Time {
	return m.CreatedAt
}

// SetCreatedAt 登録日時を設定します
func (m *FacilityMapper) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}

// GetUpdatedAt 更新日時を取得します
func (m *FacilityMapper) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

// SetUpdatedAt 更新日時を設定します
func (m *FacilityMapper) SetUpdatedAt(t time.Time) {
	m.UpdatedAt = t
}
