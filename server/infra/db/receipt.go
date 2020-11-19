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

	err := r.sess.PutResource(newIRMapper(*m.IR))
	if err != nil {
		return errors.WithStack(err)
	}

	for _, item := range m.ReceiptItems {
		if err := r.sess.PutResource(newREMapper(*item.RE)); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func newIRMapper(m model.IR) *IRMapper {
	return &IRMapper{
		IR: m,
	}
}

// IRMapper IRモデルのリソースへのマッパー構造体
type IRMapper struct {
	model.IR
	ID        string    `dynamo:"ID,hash"`
	Metadata  string    `dynamo:"Metadata,range"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
}

// GetID IDの取得
func (m *IRMapper) GetID() string {
	return fmt.Sprintf("%s#%s", m.FacilityID, m.InvoiceYM)
}

// SetID IDの 設定
func (m *IRMapper) SetID() {
	m.ID = m.GetID()
}

// GetMetadata Metadataの取得
func (m *IRMapper) GetMetadata() string {
	return fmt.Sprintf("%d", m.Payer)
}

// SetMetadata Metadataの設定
func (m *IRMapper) SetMetadata() {
	m.Metadata = m.GetMetadata()
}

// SetCreatedAt 登録日時の設定
func (m *IRMapper) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}

func newREMapper(m model.RE) *REMapper {
	return &REMapper{
		RE: m,
	}
}

// REMapper REerモデルのリソースへのマッパー構造体
type REMapper struct {
	model.RE
	ID        string    `dynamo:"ID,hash"`
	Metadata  string    `dynamo:"Metadata,range"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
}

// GetID IDの取得
func (m *REMapper) GetID() string {
	return fmt.Sprintf("%s#%s", m.GetFacilityID(), m.GetInvoiceYM())
}

// SetID IDの 設定
func (m *REMapper) SetID() {
	m.ID = m.GetID()
}

// GetMetadata Metadataの取得
func (m *REMapper) GetMetadata() string {
	return fmt.Sprintf("%d#%d", m.GetReceiptNo(), m.GetIndex())
}

// SetMetadata Metadataの設定
func (m *REMapper) SetMetadata() {
	m.Metadata = m.GetMetadata()
}

// SetCreatedAt 登録日時の設定
func (m *REMapper) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}
