package model

// DailyClientPoints 日別患者行為点数の集合
type DailyClientPoints []DailyClientPoint

// Point 行為点数
type Point struct {
	// 施設コード
	FacilityID string
	// 病棟コード
	BuildingID string
	// 入院基本料区分
	PlanClass string
	// 行為点数
	Value int64
}

// DailyClientPoint 日別患者行為点数
type DailyClientPoint struct {
	// 実施年月日
	CaredOn string
	// データ識別番号
	ClientID string
	// 行為点数
	Point
}
