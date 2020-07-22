package model

// Facilities 施設モデルの集合
type Facilities []Facility

// Facility 施設モデル
type Facility struct {
	// 施設コード
	FacilityID string
	// 施設名
	FacilityName string
}
