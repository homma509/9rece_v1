package model

// Receipt レセプト
type Receipt struct {
	IR           *IR                     // 医療機関情報レコード
	ReceiptItems map[uint32]*ReceiptItem // レセプト明細
}

// Container レセプト番号指定のレセプト明細の取得
func (r *Receipt) Container(key uint32) *ReceiptItem {
	if r.ReceiptItems == nil {
		r.ReceiptItems = map[uint32]*ReceiptItem{}
	}
	if _, ok := r.ReceiptItems[key]; !ok {
		r.ReceiptItems[key] = &ReceiptItem{}
	}
	return r.ReceiptItems[key]
}

// ReceiptItem レセプト明細
type ReceiptItem struct {
	RE  *RE   // レセプト共通レコード
	SYs []*SY // 傷病名レコード
	SIs []*SI // 診療行為レコード
	COs []*CO // コメントコード
}
