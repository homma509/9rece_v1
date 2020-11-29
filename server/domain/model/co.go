package model

const (
	// CORecordType コメントレコードのレコード識別情報
	CORecordType = "CO"
)

// CO コメントレコード
type CO struct {
	RecordType    string // レコード識別情報
	TreatmentType string // 診療識別
	ChargeType    string // 負担区分
	CommentID     string // コメントコード
	Comment       string // 文字データ
}
