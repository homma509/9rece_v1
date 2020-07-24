package usecase

import (
	"context"

	"github.com/homma509/9rece/backend/domain/model"
	"github.com/homma509/9rece/backend/domain/repository"
)

// DailyClientPointUsecase 日別患者行為点数ユースケースのインターフェース
type DailyClientPointUsecase interface {
	Store(context.Context, model.DailyClientPoints) error
}

type dailyClientPointUsecase struct {
	dailyClientPointRepository repository.DailyClientPointRepository
}

// NewDailyClientPointUsecase 日別患者別行為点数ユースケースを生成します
func NewDailyClientPointUsecase(r repository.DailyClientPointRepository) DailyClientPointUsecase {
	return &dailyClientPointUsecase{
		dailyClientPointRepository: r,
	}
}

// Store 日別患者別行為点数を登録します
func (u *dailyClientPointUsecase) Store(ctx context.Context, points model.DailyClientPoints) error {
	return u.dailyClientPointRepository.Save(ctx, points)
}
