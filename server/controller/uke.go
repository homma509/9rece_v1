package controller

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/homma509/9rece/server/domain/model"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"golang.org/x/xerrors"
)

// UkeController UKEコントローラのインターフェース
type UkeController interface {
	Move(context.Context, events.S3Event) error
}

// UkeFile UKEファイルのインターフェース
type UkeFile interface {
	GetObject(bucket, key string) (io.ReadCloser, error)
	MoveObject(srcBucket, srcKey, dstBucket, dstKey string) error
}

type ukeController struct {
	ukeFile      UkeFile
	serverBucket string
}

// NewUkeController UKEコントローラを生成します
func NewUkeController(f UkeFile, serverBucket string) UkeController {
	return &ukeController{
		ukeFile:      f,
		serverBucket: serverBucket,
	}
}

// Move UKEファイルを移動します
func (c *ukeController) Move(ctx context.Context, event events.S3Event) error {
	for _, record := range event.Records {
		bucket, _ := url.QueryUnescape(record.S3.Bucket.Name)
		key, _ := url.QueryUnescape(record.S3.Object.Key)

		err := c.move(ctx, bucket, key)
		if err != nil {
			return xerrors.Errorf("on Move bucket %s key %s: %w", bucket, key, err)
		}
	}
	return nil
}

func (c *ukeController) move(ctx context.Context, bucket, key string) error {
	// UKEファイルの読込
	f, err := c.ukeFile.GetObject(bucket, key)
	if err != nil {
		return xerrors.Errorf("on move.GetObject: %w", err)
	}
	defer f.Close()

	// 移動先パスの取得
	path, err := c.path(f)
	if err != nil {
		return xerrors.Errorf("on move.path: %w", err)
	}

	// UKEファイルの移動
	if err = c.ukeFile.MoveObject(bucket, key, c.serverBucket, path); err != nil {
		return xerrors.Errorf("on move.MoveObject: %w", err)
	}

	return nil
}

func (c *ukeController) path(f io.ReadCloser) (string, error) {
	r := csv.NewReader(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))
	r.Comma = ','
	r.FieldsPerRecord = -1
	r.ReuseRecord = true

	record, err := r.Read()
	if err == io.EOF {
		return "", xerrors.Errorf("on path.Read uke file EOF: %w", err)
	}
	if err != nil {
		return "", xerrors.Errorf("on path.Read uke file empty: %w", err)
	}

	payer, err := strconv.ParseUint(record[1], 10, 8)
	if err != nil {
		return "", xerrors.Errorf("on path.ParseUnit Payer couldn't convert number from %v: %w", record[1], err)
	}
	prefecture, err := strconv.ParseUint(record[2], 10, 8)
	if err != nil {
		return "", xerrors.Errorf("on path.ParseUnit Prefecture couldn't convert number from %v: %w", record[2], err)
	}
	pointTable, err := strconv.ParseUint(record[3], 10, 8)
	if err != nil {
		return "", xerrors.Errorf("on path.ParseUnit PointTable couldn't convert number from %v: %w", record[3], err)
	}
	medicalNo, err := strconv.ParseUint(record[4], 10, 32)
	if err != nil {
		return "", xerrors.Errorf("on path.ParseUnit MedicalNo couldn't convert number from %v: %w", record[4], err)
	}
	invoiceYearMonth, err := strconv.ParseUint(record[7], 10, 32)
	if err != nil {
		return "", xerrors.Errorf("on path.ParseUnit InvoiceYearMonth couldn't convert number from %v: %w", record[7], err)
	}
	multiVolumeID, err := strconv.ParseUint(record[8], 10, 8)
	if err != nil {
		return "", xerrors.Errorf("on path.ParseUnit MultiVolumeID couldn't convert number from %v: %w", record[8], err)
	}

	var ir model.IR
	ir.RecordID = record[0]
	ir.Payer = uint8(payer)
	ir.Prefecture = uint8(prefecture)
	ir.PointTable = uint8(pointTable)
	ir.MedicalFacilityNo = uint32(medicalNo)
	ir.Reserve = record[5]
	ir.MedicalFacilityName = record[6]
	ir.InvoiceYearMonth = uint32(invoiceYearMonth)
	ir.MultiVolumeID = uint8(multiVolumeID)
	ir.Phone = record[9]

	return fmt.Sprintf("uke/%d/%d_%s_%d.UKE",
		ir.MedicalFacilityNo,
		ir.MedicalFacilityNo,
		ir.MedicalFacilityName,
		ir.InvoiceYearMonth,
	), nil
}
