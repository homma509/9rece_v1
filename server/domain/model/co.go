package model

const (
	// CORecordType コメントレコードのレコード識別情報
	CORecordType = "CO"
)

// CO コメントレコード
type CO struct {
	FacilityID string `dynamo:"-"` // 医療機関コード
	InvoiceYM  string `dynamo:"-"` // 請求年月
	Index      uint64 `dynamo:"-"` // インデックス
	ReceiptNo  uint32 `dynamo:"-"` // レセプト番号

	RecordType    string // レコード識別情報
	TreatmentType uint8  // 診療識別
	ChargeType    string // 負担区分
	CommentID     string // コメントコード
	Comment       string // 文字データ
}

// GetFacilityID 医療機関コード
func (co *CO) GetFacilityID() string {
	return co.FacilityID
}

// GetInvoiceYM 請求年月
func (co *CO) GetInvoiceYM() string {
	return co.InvoiceYM
}

// GetIndex インデックス
func (co *CO) GetIndex() uint64 {
	return co.Index
}

// GetReceiptNo レセプト番号
func (co *CO) GetReceiptNo() uint32 {
	return co.ReceiptNo
}
