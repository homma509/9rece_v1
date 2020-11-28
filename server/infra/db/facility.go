package db

import (
	"context"
	"time"

	"github.com/homma509/9rece/server/domain/model"
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
	// rs := []Resource{}
	// for _, p := range ps {
	// 	rs = append(rs, r.newFacilityMapper(p))
	// }

	// err := r.sess.PutResources(rs)
	// if err != nil {
	// 	return errors.WithStack(err)
	// }

	return nil
}

func (r *FacilityRepository) newFacilityMapper(p model.Facility) *FacilityMapper {
	new := &FacilityMapper{
		Facility: p,
	}

	// var old FacilityMapper
	// err := r.sess.GetResource(new, &old)
	// if err == nil {
	// 	new.SetPK()
	// 	new.SetSK()
	// 	new.SetCreatedAt(old.CreatedAt)
	// 	new.SetVersion(old.Version)
	// }

	return new
}

// FacilityMapper 施設モデルのリソースへのマッパー構造体
type FacilityMapper struct {
	model.Facility
	ID        string    `dynamo:"ID,hash"`
	Metadata  string    `dynamo:"Metadata,range"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
}

// GetID IDの取得
func (m *FacilityMapper) GetID() string {
	return ""
}

// SetID IDの 設定
func (m *FacilityMapper) SetID() {
	m.ID = m.GetID()
}

// GetMetadata Metadataの取得
func (m *FacilityMapper) GetMetadata() string {
	return ""
}

// SetMetadata Metadataの設定
func (m *FacilityMapper) SetMetadata() {
	m.Metadata = m.GetMetadata()
}

// SetCreatedAt 登録日時の設定
func (m *FacilityMapper) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}
