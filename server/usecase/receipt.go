package usecase

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/homma509/9rece/server/domain/model"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"golang.org/x/xerrors"
)

// ReceiptUsecase レセプトユースケースのインターフェース
type ReceiptUsecase interface {
	Move(ctx context.Context, bucket, key string) error
	// Path(context.Context, file.File) (string, error)
	// Store(context.Context, model.Receipt) error
}

// ReceiptFile レセプトファイルのインターフェース
type ReceiptFile interface {
	GetObject(bucket, key string) (io.ReadCloser, error)
	MoveObject(srcBucket, srcKey, dstBucket, dstKey string) error
}

type receiptUsecase struct {
	dstBucket string
	file      ReceiptFile
	// receiptRepository repository.ReceiptRepository
}

// NewReceiptUsecase レセプトユースケースの生成
func NewReceiptUsecase(f ReceiptFile, serverBucket string) ReceiptUsecase {
	return &receiptUsecase{
		file:      f,
		dstBucket: serverBucket,
	}
}

// // Store レセプトの登録
// func (u *receiptUsecase) Store(ctx context.Context, receipt model.Receipt) error {
// 	return u.receiptRepository.Save(ctx, receipt)
// }

func (u *receiptUsecase) Move(ctx context.Context, bucket, key string) error {
	// レセプトファイルの読込
	f, err := u.file.GetObject(bucket, key)
	if err != nil {
		return xerrors.Errorf("on Move.GetObject: %w", err)
	}
	defer f.Close()

	// 移動先パスの取得
	ir, err := readIR(f)
	if err != nil {
		return xerrors.Errorf("on Move.readIR: %w", err)
	}
	path := path(ir)

	// レセプトファイルの移動
	if err = u.file.MoveObject(bucket, key, u.dstBucket, path); err != nil {
		return xerrors.Errorf("on Move.MoveObject: %w", err)
	}

	return nil
}

func readIR(f io.ReadCloser) (*model.IR, error) {
	r := csv.NewReader(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))
	r.Comma = ','
	r.FieldsPerRecord = -1
	r.ReuseRecord = true

	record, err := r.Read()
	if err == io.EOF {
		return nil, xerrors.Errorf("on readIR.Read receipt file EOF: %w", err)
	}
	if err != nil {
		return nil, xerrors.Errorf("on readIR.Read receipt file empty: %w", err)
	}

	ir, err := ir(record)
	if err != nil {
		return nil, xerrors.Errorf("on readIR.Read receipt file empty: %w", err)
	}

	return ir, nil
}

func ir(record []string) (*model.IR, error) {
	payer, err := strconv.ParseUint(record[1], 10, 8)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Payer couldn't convert number from %v: %w", record[1], err)
	}
	pointTable, err := strconv.ParseUint(record[3], 10, 8)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit PointTable couldn't convert number from %v: %w", record[3], err)
	}

	return &model.IR{
		RecordType:    record[0],
		Payer:         uint8(payer),
		Prefecture:    record[2],
		PointTable:    uint8(pointTable),
		FacilityID:    record[4],
		Reserve:       record[5],
		FacilityName:  record[6],
		InvoiceYM:     record[7],
		MultiVolumeNo: record[8],
		Phone:         record[9],
	}, nil
}

func path(ir *model.IR) string {
	return fmt.Sprintf("receipts/%d/%d_%s_%d.UKE",
		ir.FacilityID,
		ir.FacilityID,
		ir.FacilityName,
		ir.InvoiceYM,
	)
}
