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
	CreatedAt time.Time `dynamo:"CreatedAt"`
	UpdatedAt time.Time `dynamo:"UpdatedAt"`
}

// EntityName Entity名の取得
func (r *FacilityMapper) EntityName() string {
	t := reflect.TypeOf(r.Facility)
	return t.Name()
}

// GetPK PKを取得します
func (r *FacilityMapper) GetPK() string {
	return r.FacilityID
}

// SetPK PKを設定します
func (r *FacilityMapper) SetPK() {
	r.PK = r.GetPK()
}

// GetSK SKを取得します
func (r *FacilityMapper) GetSK() string {
	return fmt.Sprintf("%sInfo", r.EntityName())
}

// SetSK SKを設定します
func (r *FacilityMapper) SetSK() {
	r.SK = r.GetSK()
}

// GetVersion バージョンを取得します
func (r *FacilityMapper) GetVersion() uint64 {
	return r.Version
}

// SetVersion Versionを設定します
func (r *FacilityMapper) SetVersion(v uint64) {
	r.Version = v
}

// GetCreatedAt 登録日時を取得します
func (r *FacilityMapper) GetCreatedAt() time.Time {
	return r.CreatedAt
}

// SetCreatedAt 登録日時を設定します
func (r *FacilityMapper) SetCreatedAt(t time.Time) {
	r.CreatedAt = t
}

// GetUpdatedAt 更新日時を取得します
func (r *FacilityMapper) GetUpdatedAt() time.Time {
	return r.UpdatedAt
}

// SetUpdatedAt 更新日時を設定します
func (r *FacilityMapper) SetUpdatedAt(t time.Time) {
	r.UpdatedAt = t
}
