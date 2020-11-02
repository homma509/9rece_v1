package usecase

import (
	"context"

	"github.com/homma509/9rece/backend/domain/model"
)

// URLUsecase URLユースケースのインターフェース
type URLUsecase interface {
	Get(context.Context) (*model.URL, error)
}

type urlUsecase struct {
}

// NewURLUsecase URLユースケースを生成します
func NewURLUsecase() URLUsecase {
	return &urlUsecase{}
}

// Get URLを取得します
func (u *urlUsecase) Get(ctx context.Context) (*model.URL, error) {
	return &model.URL{
		URL: "https://www.golang.org",
	}, nil
}
