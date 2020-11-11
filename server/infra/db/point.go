package db

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/homma509/9rece/server/domain/model"
	"github.com/pkg/errors"
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
	rs := []Resource{}
	for _, p := range ps {
		rs = append(rs, r.newDailyClientPointMapper(p))
	}

	err := r.sess.PutResources(rs)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *DailyClientPointRepository) newDailyClientPointMapper(p model.DailyClientPoint) *DailyClientPointMapper {
	new := &DailyClientPointMapper{
		DailyClientPoint: p,
	}

	var old DailyClientPointMapper
	err := r.sess.GetResource(new, &old)
	if err == nil {
		new.SetPK()
		new.SetSK()
		new.SetCreatedAt(old.CreatedAt)
		new.SetVersion(old.Version)
	}

	return new
}

// DailyClientPointMapper 日別患者行為点数モデルのリソースへのマッパー構造体
type DailyClientPointMapper struct {
	model.DailyClientPoint
	PK        string    `dynamo:"ID,hash"`
	SK        string    `dynamo:"DataType,range"`
	Version   uint64    `dynamo:"Version"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
	UpdatedAt time.Time `dynamo:"UpdatedAt"`
}

// EntityName Entity名の取得
func (m *DailyClientPointMapper) EntityName() string {
	t := reflect.TypeOf(m.DailyClientPoint)
	return t.Name()
}

// GetPK PKを取得します
func (m *DailyClientPointMapper) GetPK() string {
	return fmt.Sprintf("%s:%s", m.FacilityID, m.CaredOn[:7])
}

// SetPK PKを設定します
func (m *DailyClientPointMapper) SetPK() {
	m.PK = m.GetPK()
}

// GetSK SKを取得します
func (m *DailyClientPointMapper) GetSK() string {
	return fmt.Sprintf("%sInfo_#%s:%s", m.EntityName(), m.ClientID, m.CaredOn)
}

// SetSK SKを設定します
func (m *DailyClientPointMapper) SetSK() {
	m.SK = m.GetSK()
}

// GetVersion バージョンを取得します
func (m *DailyClientPointMapper) GetVersion() uint64 {
	return m.Version
}

// SetVersion Versionを設定します
func (m *DailyClientPointMapper) SetVersion(v uint64) {
	m.Version = v
}

// GetCreatedAt 登録日時を取得します
func (m *DailyClientPointMapper) GetCreatedAt() time.Time {
	return m.CreatedAt
}

// SetCreatedAt 登録日時を設定します
func (m *DailyClientPointMapper) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}

// GetUpdatedAt 更新日時を取得します
func (m *DailyClientPointMapper) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

// SetUpdatedAt 更新日時を設定します
func (m *DailyClientPointMapper) SetUpdatedAt(t time.Time) {
	m.UpdatedAt = t
}
