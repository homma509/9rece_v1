package model

const (
	// IRRecordType 医療機関情報レコードのレコード識別情報
	IRRecordType = "IR"
)

// IR 医療機関情報レコード
type IR struct {
	RecordType    string // レコード識別情報
	Payer         uint8  // 審査支払機関
	Prefecture    string // 都道府県
	PointTable    uint8  // 点数表
	FacilityID    string // 医療機関コード
	Reserve       string // 予備
	FacilityName  string // 医療機関名称
	InvoiceYM     string // 請求年月
	MultiVolumeNo string // マルチボリューム識別情報
	Phone         string // 電話番号
}
