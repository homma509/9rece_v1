package db

import (
	"context"
	"time"

	"github.com/homma509/9rece/server/domain/model"
)

// DailyClientPointRepository 日別患者行為点数リポジトリの構造体
type DailyClientPointRepository struct {
	sess *Session
}

// NewDailyClientPointRepository 日別患者行為点数リポジトリを生成します
func NewDailyClientPointRepository(sess *Session) *DailyClientPointRepository {
	return &DailyClientPointRepository{
		sess: sess,
	}
}

// Save 日別患者行為点数のスライスを登録します
func (r *DailyClientPointRepository) Save(ctx context.Context, ps model.DailyClientPoints) error {
	// rs := []Resource{}
	// for _, p := range ps {
	// 	rs = append(rs, r.newDailyClientPointMapper(p))
	// }

	// err := r.sess.PutResources(rs)
	// if err != nil {
	// 	return errors.WithStack(err)
	// }

	return nil
}

func (r *DailyClientPointRepository) newDailyClientPointMapper(p model.DailyClientPoint) *DailyClientPointMapper {
	new := &DailyClientPointMapper{
		DailyClientPoint: p,
	}

	// var old DailyClientPointMapper
	// err := r.sess.GetResource(new, &old)
	// if err == nil {
	// 	new.SetPK()
	// 	new.SetSK()
	// 	new.SetCreatedAt(old.CreatedAt)
	// 	new.SetVersion(old.Version)
	// }

	return new
}

// DailyClientPointMapper 日別患者行為点数モデルのリソースへのマッパー構造体
type DailyClientPointMapper struct {
	model.DailyClientPoint
	ID        string    `dynamo:"ID,hash"`
	Metadata  string    `dynamo:"Metadata,range"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
}

// GetID IDの取得
func (m *DailyClientPointMapper) GetID() string {
	return ""
}

// SetID IDの 設定
func (m *DailyClientPointMapper) SetID() {
	m.ID = m.GetID()
}

// GetMetadata Metadataの取得
func (m *DailyClientPointMapper) GetMetadata() string {
	return ""
}

// SetMetadata Metadataの設定
func (m *DailyClientPointMapper) SetMetadata() {
	m.Metadata = m.GetMetadata()
}

// SetCreatedAt 登録日時の設定
func (m *DailyClientPointMapper) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}
