package usecase

import (
	"context"

	"github.com/homma509/9rece/backend/domain/model"
	"github.com/homma509/9rece/backend/domain/repository"
)

// FacilityUsecase 施設ユースケースのインターフェース
type FacilityUsecase interface {
	Store(context.Context, model.Facilities) error
}

type facilityUsecase struct {
	facilityRepository repository.FacilityRepository
}

// NewFacilityUsecase 施設ユースケースを生成します
func NewFacilityUsecase(fr repository.FacilityRepository) FacilityUsecase {
	return &facilityUsecase{
		facilityRepository: fr,
	}
}

// Store 施設を登録します
func (f *facilityUsecase) Store(ctx context.Context, facilities model.Facilities) error {
	return f.facilityRepository.Save(ctx, facilities)
}
