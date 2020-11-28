package controller

import (
	"context"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/homma509/9rece/server/usecase"
	"golang.org/x/xerrors"
)

// ReceiptController レセプトコントローラのインターフェース
type ReceiptController interface {
	Move(context.Context, events.S3Event) error
	Store(context.Context, events.S3Event) error
}

type receiptController struct {
	usecase usecase.ReceiptUsecase
}

// NewReceiptController レセプトコントローラの生成
func NewReceiptController(u usecase.ReceiptUsecase) ReceiptController {
	return &receiptController{
		usecase: u,
	}
}

// Move レセプトファイルの移動
func (c *receiptController) Move(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {
		bucket, _ := url.QueryUnescape(record.S3.Bucket.Name)
		key, _ := url.QueryUnescape(record.S3.Object.Key)

		err := c.usecase.Move(ctx, bucket, key)
		if err != nil {
			return xerrors.Errorf("on Move bucket %s key %s: %w", bucket, key, err)
		}
	}
	return nil
}

// Store レセプトファイルの登録
func (c *receiptController) Store(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {
		bucket, _ := url.QueryUnescape(record.S3.Bucket.Name)
		key, _ := url.QueryUnescape(record.S3.Object.Key)

		err := c.usecase.Store(ctx, bucket, key)
		if err != nil {
			return xerrors.Errorf("on Store bucket %s key %s: %w", bucket, key, err)
		}
	}
	return nil
}
