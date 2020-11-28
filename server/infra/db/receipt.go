package db

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/homma509/9rece/server/domain/model"
	"github.com/pkg/errors"
)

// GeentityNameFromStruct StructからEntity名を取得
func entityNameFromStruct(s interface{}) string {
	return reflect.TypeOf(s).Name()
}

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
func (r *ReceiptRepository) Save(ctx context.Context, m model.Receipt) error {
	// TODO 登録前にIR単位で全削除を実施し、冪等にする
	rs := []Resource{}
	rs = append(rs, newIRMapper(m.IR))

	for _, item := range m.ReceiptItems() {
		rs = append(rs, newREMapper(m.IR, item.RE))
		for i, sy := range item.SYs {
			rs = append(rs, newSYMapper(m.IR, item.RE, sy, i))
		}
	}

	err := r.sess.PutResources(rs)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func newIRMapper(ir model.IR) *IRMapper {
	return &IRMapper{
		IR: ir,
	}
}

// IRMapper IRモデルのリソースへのマッパー構造体
type IRMapper struct {
	model.IR
	ID        string    `dynamo:"ID,hash"`
	Metadata  string    `dynamo:"Metadata,range"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
}

// SetID IDの 設定
func (m *IRMapper) SetID() {
	m.ID = fmt.Sprintf("%s#%s", m.FacilityID, m.InvoiceYM)
}

// SetMetadata Metadataの設定
func (m *IRMapper) SetMetadata() {
	m.Metadata = fmt.Sprintf("%d#%s", m.Payer, entityNameFromStruct(m.IR))
}

// SetCreatedAt 登録日時の設定
func (m *IRMapper) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}

func newREMapper(ir model.IR, re model.RE) *REMapper {
	return &REMapper{
		RE:         re,
		FacilityID: ir.FacilityID,
		InvoiceYM:  ir.InvoiceYM,
	}
}

// REMapper REモデルのリソースへのマッパー構造体
type REMapper struct {
	model.RE
	ID         string    `dynamo:"ID,hash"`
	Metadata   string    `dynamo:"Metadata,range"`
	CreatedAt  time.Time `dynamo:"CreatedAt"`
	FacilityID string    `dynamo:"FacilityID"`
	InvoiceYM  string    `dynamo:"InvoiceYM"`
}

// SetID IDの 設定
func (m *REMapper) SetID() {
	m.ID = fmt.Sprintf("%s#%s", m.FacilityID, m.InvoiceYM)
}

// SetMetadata Metadataの設定
func (m *REMapper) SetMetadata() {
	m.Metadata = fmt.Sprintf("%d#%s", m.ReceiptNo, entityNameFromStruct(m.RE))
}

// SetCreatedAt 登録日時の設定
func (m *REMapper) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}

func newSYMapper(ir model.IR, re model.RE, sy model.SY, i int) *SYMapper {
	return &SYMapper{
		SY:         sy,
		FacilityID: ir.FacilityID,
		InvoiceYM:  ir.InvoiceYM,
		ReceiptNo:  re.ReceiptNo,
		Index:      uint64(i),
	}
}

// SYMapper SYモデルのリソースへのマッパー構造体
type SYMapper struct {
	model.SY
	ID         string    `dynamo:"ID,hash"`
	Metadata   string    `dynamo:"Metadata,range"`
	CreatedAt  time.Time `dynamo:"CreatedAt"`
	FacilityID string    `dynamo:"FacilityID"`
	InvoiceYM  string    `dynamo:"InvoiceYM"`
	ReceiptNo  uint32    `dynamo:"ReceiptNo"`
	Index      uint64    `dynamo:"Index"`
}

// SetID IDの 設定
func (m *SYMapper) SetID() {
	m.ID = fmt.Sprintf("%s#%s", m.FacilityID, m.InvoiceYM)
}

// SetMetadata Metadataの設定
func (m *SYMapper) SetMetadata() {
	m.Metadata = fmt.Sprintf("%d#%s#%d", m.ReceiptNo, entityNameFromStruct(m.SY), m.Index)
}

// SetCreatedAt 登録日時の設定
func (m *SYMapper) SetCreatedAt(t time.Time) {
	m.CreatedAt = t
}
