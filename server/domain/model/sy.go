package model

// SY 傷病名レコード
type SY struct {
	RecordType  string // レコード識別情報
	DiseaseID   string // 傷病名コード
	ReceiptedAt string // 診療開始日
	OutcomeType uint8  // 転記区分
	ModifierID  string // 修飾語コード
	DiseaseName string // 傷病名称
	MainDisease uint8  // 主傷病
	Comment     string // 補足コメント
}
