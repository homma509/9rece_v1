package model

// CO コメントレコード
type CO struct {
	RecordType    string // レコード識別情報
	TreatmentType uint8  // 診療識別
	ChargeType    string // 負担区分
	CommentID     string // コメントコード
	Comment       string // 文字データ
}
