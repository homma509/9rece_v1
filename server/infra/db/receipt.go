package db

import (
	"context"
	"fmt"
	"time"

	"github.com/homma509/9rece/server/domain/model"
	"github.com/pkg/errors"
)

// ReceiptRepository レセプトリポジトリの構造体
type ReceiptRepository struct {
	sess *Session
}

// NewReceiptRepository レセプトリポジトリの生成
func NewReceiptRepository(sess *Session) *ReceiptRepository {
	return &ReceiptRepository{
		sess: sess,
	}
}

// Save レセプトの登録
func (r *ReceiptRepository) Save(ctx context.Context, m *model.Receipt) error {
	// TODO 登録前に全削除を実施し、冪等にする
	err := r.sess.PutResource(r.newReceiptMapper(m))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *ReceiptRepository) newReceiptMapper(m *model.Receipt) *ReceiptMapper {
	return &ReceiptMapper{
		IR: *m.IR,
	}
}

// ReceiptMapper レセプトモデルのリソースへのマッパー構造体
type ReceiptMapper struct {
	model.IR
	ID        string    `dynamo:"ID,hash"`
	Metadata  string    `dynamo:"Metadata,range"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
}

// GetID IDの取得
func (m *ReceiptMapper) GetID() string {
	return fmt.Sprintf("%s#%s", m.FacilityID, m.InvoiceYM)
}

// SetID IDの 設定
func (m *ReceiptMapper) SetID() {
	m.ID = m.GetID()
}

// GetMetadata Metadataの取得
func (m *ReceiptMapper) GetMetadata() string {
	return fmt.Sprintf("%d", m.Payer)
}

// SetMetadata Metadataの設定
func (m *ReceiptMapper) SetMetadata() {
	m.Metadata = m.GetMetadata()
}

// SetCreatedAt 登録日時の設定
func (m *ReceiptMapper) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}
