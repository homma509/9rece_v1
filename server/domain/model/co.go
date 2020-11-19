package model

const (
	// CORecordType コメントレコードのレコード識別情報
	CORecordType = "CO"
)

// CO コメントレコード
type CO struct {
	FacilityID    string // 医療機関コード
	InvoiceYM     string // 請求年月
	Index         uint64 // インデックス
	ReceiptNo     uint32 // レセプト番号
	RecordType    string // レコード識別情報
	TreatmentType uint8  // 診療識別
	ChargeType    string // 負担区分
	CommentID     string // コメントコード
	Comment       string // 文字データ
}
