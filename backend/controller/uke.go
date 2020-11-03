package controller

import (
	"context"
	"io"

	"github.com/aws/aws-lambda-go/events"
)

// UkeController UKEコントローラのインターフェース
type UkeController interface {
	Move(context.Context, events.S3Event) error
}

// UkeFile UKEファイルのインターフェース
type UkeFile interface {
	GetObject(bucket, key string) (io.ReadCloser, error)
	MoveObject(bucket, src, dst string) error
}

type ukeController struct {
	ukeFile UkeFile
}

// NewUkeController UKEコントローラを生成します
func NewUkeController(f UkeFile) UkeController {
	return &ukeController{
		ukeFile: f,
	}
}

// Move UKEファイルを移動します
func (c *ukeController) Move(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		err := c.move(ctx, bucket, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ukeController) move(ctx context.Context, bucket, key string) error {
	// // S3ファイル取得
	// f, err := c.ukeFile.GetObject(bucket, key)
	// if err != nil {
	// 	return err
	// }
	// defer f.Close()

	// // ファイル読み込み
	// facilities, err := c.read(f)
	// if err != nil {
	// 	return err
	// }

	// S3ファイルコピー
	err := c.ukeFile.MoveObject(bucket, key, "uke/テスト施設_202011.uke")
	if err != nil {
		return err
	}

	return nil
}

// func (c *ukeController) read(f io.ReadCloser) (model.Facilities, error) {
// 	facilities := model.Facilities{}

// 	r := csv.NewReader(f)
// 	r.Comma = '\t'
// 	r.FieldsPerRecord = 2
// 	r.ReuseRecord = true

// 	for {
// 		record, err := r.Read()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return nil, err
// 		}

// 		ukeID := record[0]
// 		_, err = strconv.ParseInt(record[0], 10, 64)
// 		if err != nil {
// 			// UKEコードが数値以外はHeaderとみなす
// 			continue
// 		}
// 		ukeName := record[1]

// 		uke := model.Uke{
// 			UkeID:   ukeID,
// 			UkeName: ukeName,
// 		}

// 		facilities = append(facilities, uke)
// 	}

// 	return facilities, nil
// }
