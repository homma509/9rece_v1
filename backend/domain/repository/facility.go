package repository

import (
	"context"

	"github.com/homma509/9rece/backend/domain/model"
)

// FacilityRepository 施設モデルのリポジトリのインターフェース
type FacilityRepository interface {
	Save(context.Context, model.Facilities) error
}
