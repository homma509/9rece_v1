package model

// IR 医療機関情報レコード
type IR struct {
	RecordID            string // レコード識別情報
	Payer               uint8  // 審査支払機関
	Prefecture          uint8  // 都道府県
	PointTable          uint8  // 点数表
	MedicalFacilityNo   uint32 // 医療機関コード
	Reserve             string // 予備
	MedicalFacilityName string // 医療機関名称
	InvoiceYearMonth    uint32 // 請求年月
	MultiVolumeID       uint8  // マルチボリューム識別情報
	Phone               string // 電話番号
}
