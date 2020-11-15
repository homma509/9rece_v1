package usecase

// import (
// 	"context"

// 	"github.com/homma509/9rece/server/domain/model"
// 	"github.com/homma509/9rece/server/domain/repository"
// )

// // UkeUsecase UKEユースケースのインターフェース
// type UkeUsecase interface {
// 	Store(context.Context, model.Receipt) error
// }

// type ukeUsecase struct {
// 	ukeRepository repository.UkeRepository
// }

// // NewUkeUsecase UKEユースケースの生成
// func NewUkeUsecase(r repository.UkeRepository) UkeUsecase {
// 	return &ukeUsecase{
// 		ukeRepository: r,
// 	}
// }

// // Store UKEの登録
// func (u *ukeUsecase) Store(ctx context.Context, facilities model.Facilities) error {
// 	return u.facilityRepository.Save(ctx, facilities)
// }
