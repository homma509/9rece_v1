package model

// Receipt レセプト
type Receipt struct {
	IR           *IR           // 医療機関情報
	ReceiptItems []ReceiptItem // レセプト明細
}

// ReceiptItem レセプト明細
type ReceiptItem struct {
	RE  *RE   // レセプト共通レコード
	SYs []*SY // 傷病名レコード
	SIs []*SI // 診療行為レコード
	COs []*CO // コメントコード
}
