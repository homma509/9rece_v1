package controller

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/homma509/9rece/backend/domain/model"
	"github.com/homma509/9rece/backend/usecase"
)

// DailyClientPointController 日別患者行為点数コントローラのインターフェース
type DailyClientPointController interface {
	Post(context.Context, events.S3Event) error
}

// DailyClientPointFile  日別患者行為点数ファイルのインターフェース
type DailyClientPointFile interface {
	GetObject(string, string) (io.ReadCloser, error)
}

type dailyClientPointController struct {
	dailyClientPointUsecase usecase.DailyClientPointUsecase
	dailyClientPointFile    DailyClientPointFile
}

// NewDailyClientPointController 日別患者行為点数コントローラを生成します
func NewDailyClientPointController(u usecase.DailyClientPointUsecase, f DailyClientPointFile) DailyClientPointController {
	return &dailyClientPointController{
		dailyClientPointUsecase: u,
		dailyClientPointFile:    f,
	}
}

// Post 日別患者行為点数を登録します
func (c *dailyClientPointController) Post(ctx context.Context, event events.S3Event) error {
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

func (c *dailyClientPointController) store(ctx context.Context, bucket, key string) error {
	// S3ファイル取得
	f, err := c.dailyClientPointFile.GetObject(bucket, key)
	if err != nil {
		return err
	}
	defer f.Close()

	// ファイル読み込み
	points, err := c.read(f)
	if err != nil {
		return err
	}

	// DBに登録
	err = c.dailyClientPointUsecase.Store(ctx, points)
	if err != nil {
		return err
	}

	return nil
}

func (c *dailyClientPointController) read(f io.ReadCloser) (model.DailyClientPoints, error) {
	points := model.DailyClientPoints{}

	r := csv.NewReader(f)
	r.Comma = '\t'
	r.FieldsPerRecord = 6
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
		d, err := time.Parse("2006/01/02", record[1])
		if err != nil {
			log.Printf("%v couldn't parse time, %v", record[1], err)
			continue
		}
		caredOn := d.Format("2006-01-02")

		clientID := record[2]
		buildingID := record[3]
		planClass := record[4]
		value, err := strconv.ParseInt(record[5], 10, 64)
		if err != nil {
			log.Printf("%v couldn't parse int, %v", record[5], err)
			continue
		}

		point := model.DailyClientPoint{
			CaredOn:  caredOn,
			ClientID: clientID,
			Point: model.Point{
				FacilityID: facilityID,
				BuildingID: buildingID,
				PlanClass:  planClass,
				Value:      value,
			},
		}

		points = append(points, point)
	}

	return points, nil
}
