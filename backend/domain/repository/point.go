package repository

import (
	"context"

	"github.com/homma509/9rece/backend/domain/model"
)

// DailyClientPointRepository 日別患者行為点数モデルのリポジトリのインターフェース
type DailyClientPointRepository interface {
	Save(context.Context, model.DailyClientPoints) error
}
