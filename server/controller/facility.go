package controller

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/homma509/9rece/server/domain/model"
	"github.com/homma509/9rece/server/usecase"
)

// FacilityController 施設コントローラのインターフェース
type FacilityController interface {
	Post(context.Context, events.S3Event) error
}

// FacilityFile  施設ファイルのインターフェース
type FacilityFile interface {
	GetObject(string, string) (io.ReadCloser, error)
}

type facilityController struct {
	facilityUsecase usecase.FacilityUsecase
	facilityFile    FacilityFile
}

// NewFacilityController 施設コントローラを生成します
func NewFacilityController(u usecase.FacilityUsecase, f FacilityFile) FacilityController {
	return &facilityController{
		facilityUsecase: u,
		facilityFile:    f,
	}
}

// Post 施設を登録します
func (c *facilityController) Post(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		err := c.store(ctx, bucket, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *facilityController) store(ctx context.Context, bucket, key string) error {
	// S3ファイル取得
	f, err := c.facilityFile.GetObject(bucket, key)
	if err != nil {
		return err
	}
	defer f.Close()

	// ファイル読み込み
	facilities, err := c.read(f)
	if err != nil {
		return err
	}

	// DBに登録
	err = c.facilityUsecase.Store(ctx, facilities)
	if err != nil {
		return err
	}

	return nil
}

func (c *facilityController) read(f io.ReadCloser) (model.Facilities, error) {
	facilities := model.Facilities{}

	r := csv.NewReader(f)
	r.Comma = '\t'
	r.FieldsPerRecord = 2
	r.ReuseRecord = true

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		facilityID := record[0]
		_, err = strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			// 施設コードが数値以外はHeaderとみなす
			continue
		}
		facilityName := record[1]

		facility := model.Facility{
			FacilityID:   facilityID,
			FacilityName: facilityName,
		}

		facilities = append(facilities, facility)
	}

	return facilities, nil
}
