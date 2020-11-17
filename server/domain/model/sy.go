package model

const (
	// SYRecordType 傷病名レコードのレコード識別情報
	SYRecordType = "SY"
)

// SY 傷病名レコード
type SY struct {
	FacilityID string `dynamo:"-"` // 医療機関コード
	InvoiceYM  string `dynamo:"-"` // 請求年月
	Index      uint64 `dynamo:"-"` // インデックス
	ReceiptNo  uint32 `dynamo:"-"` // レセプト番号

	RecordType  string // レコード識別情報
	DiseaseID   string // 傷病名コード
	ReceiptedAt string // 診療開始日
	OutcomeType uint8  // 転記区分
	ModifierID  string // 修飾語コード
	DiseaseName string // 傷病名称
	MainDisease uint8  // 主傷病
	Comment     string // 補足コメント
}

// GetFacilityID 医療機関コード
func (sy *SY) GetFacilityID() string {
	return sy.FacilityID
}

// GetInvoiceYM 請求年月
func (sy *SY) GetInvoiceYM() string {
	return sy.InvoiceYM
}

// GetIndex インデックス
func (sy *SY) GetIndex() uint64 {
	return sy.Index
}

// GetReceiptNo レセプト番号
func (sy *SY) GetReceiptNo() uint32 {
	return sy.ReceiptNo
}
