package model

// Receipt レセプト
type Receipt struct {
	IR           IR                      // 医療機関情報レコード
	receiptItems map[uint32]*ReceiptItem // レセプト明細
}

// ReceiptItem レセプト番号指定のレセプト明細の取得
func (r *Receipt) ReceiptItem(key uint32) *ReceiptItem {
	if r.receiptItems == nil {
		r.receiptItems = map[uint32]*ReceiptItem{}
	}
	if _, ok := r.receiptItems[key]; !ok {
		r.receiptItems[key] = &ReceiptItem{
			SYs:     []SY{},
			SIInfos: []SIInfo{},
		}
	}
	return r.receiptItems[key]
}

// ReceiptItems レセプト明細リストの取得
func (r *Receipt) ReceiptItems() map[uint32]*ReceiptItem {
	return r.receiptItems
}

// ReceiptItem レセプト明細
type ReceiptItem struct {
	RE      RE       // レセプト共通レコード
	SYs     []SY     // 傷病名レコード
	SIInfos []SIInfo // 診療行為情報
}

// SIInfo 診療行為情報
type SIInfo struct {
	SI       // 診療行為レコード
	COs []CO // コメントコード
}
