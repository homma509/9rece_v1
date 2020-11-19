package usecase

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/homma509/9rece/server/domain/model"
	"github.com/homma509/9rece/server/domain/repository"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"golang.org/x/xerrors"
)

// ReceiptUsecase レセプトユースケースのインターフェース
type ReceiptUsecase interface {
	Move(ctx context.Context, bucket, key string) error
	Store(ctx context.Context, bucket, key string) error
}

// ReceiptFile レセプトファイルのインターフェース
type ReceiptFile interface {
	GetObject(bucket, key string) (io.ReadCloser, error)
	MoveObject(srcBucket, srcKey, dstBucket, dstKey string) error
}

type receiptUsecase struct {
	dstBucket string
	file      ReceiptFile
	repo      repository.ReceiptRepository
}

// NewReceiptUsecase レセプトユースケースの生成
func NewReceiptUsecase(serverBucket string, f ReceiptFile, r repository.ReceiptRepository) ReceiptUsecase {
	return &receiptUsecase{
		dstBucket: serverBucket,
		file:      f,
		repo:      r,
	}
}

// Move レセプトファイルの移動
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

// Store レセプトファイルの登録
func (u *receiptUsecase) Store(ctx context.Context, bucket, key string) error {
	// レセプトファイルの読込
	f, err := u.file.GetObject(bucket, key)
	if err != nil {
		return xerrors.Errorf("on Store.GetObject: %w", err)
	}
	defer f.Close()

	// レセプトの取得
	r, err := read(f)
	if err != nil {
		return xerrors.Errorf("on Store.readIR: %w", err)
	}

	// レセプトファイルの登録
	if err := u.repo.Save(ctx, r); err != nil {
		return xerrors.Errorf("on Store.Save: %w", err)
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

func read(f io.ReadCloser) (*model.Receipt, error) {
	r := csv.NewReader(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))
	r.Comma = ','
	r.FieldsPerRecord = -1
	r.ReuseRecord = true

	receipt := &model.Receipt{
		ReceiptItems: map[uint32]*model.ReceiptItem{},
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
			// return nil, xerrors.Errorf("on read.Read receipt file EOF: %w", err)
		}
		if err != nil {
			return nil, xerrors.Errorf("on read.Read receipt file empty: %w", err)
		}

		switch record[0] {
		case model.IRRecordType:
			ir, err := ir(record)
			if err != nil {
				return nil, xerrors.Errorf("on read.Read receipt file empty: %w", err)
			}
			receipt.IR = ir
		case model.RERecordType:
			re, err := re(record, receipt.IR)
			if err != nil {
				return nil, xerrors.Errorf("on read.Read receipt file empty: %w", err)
			}
			receipt.Container(re.ReceiptNo).RE = re
		}
	}

	return receipt, nil
}

func path(ir *model.IR) string {
	return fmt.Sprintf("receipts/%s/%s_%s_%s.UKE",
		ir.FacilityID,
		ir.FacilityID,
		ir.FacilityName,
		ir.InvoiceYM,
	)
}

func ir(record []string) (*model.IR, error) {
	if record[0] != model.IRRecordType {
		return nil, fmt.Errorf("on ir RecordType invalid value %v", record[0])
	}
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

func re(record []string, ir *model.IR) (*model.RE, error) {
	if ir == nil {
		return nil, fmt.Errorf("on re IR invalid value %v", ir)
	}
	if record[0] != model.RERecordType {
		return nil, fmt.Errorf("on re RecordType invalid value %v", record[0])
	}
	receiptNo, err := strconv.ParseUint(record[1], 10, 32)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit ReceiptType couldn't convert number from %v: %w", record[1], err)
	}
	sex, err := strconv.ParseUint(record[5], 10, 8)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Sex couldn't convert number from %v: %w", record[5], err)
	}

	return &model.RE{
		FacilityID:   ir.FacilityID,
		InvoiceYM:    ir.InvoiceYM,
		Index:        0,
		RecordType:   record[0],
		ReceiptNo:    uint32(receiptNo),
		ReceiptType:  record[2],
		ReceiptYM:    record[3],
		Name:         record[4],
		Sex:          uint8(sex),
		Birth:        record[6],
		BenefitRate:  record[7],
		AdmittedAt:   record[8],
		WardType:     record[9],
		ChargeType:   record[10],
		ReceiptNote:  record[11],
		WardCount:    record[12],
		KarteNo:      record[13],
		DiscountUnit: record[14],
		Reserve1:     record[15],
		Reserve2:     record[16],
		Reserve3:     record[17],
		SearchNo:     record[18],
		Reserve4:     record[19],
		InvoiceInfo:  record[20],
		Subject1:     record[21],
		Part1:        record[22],
		Sex1:         record[23],
		Treatment1:   record[24],
		Disease1:     record[25],
		Subject2:     record[26],
		Part2:        record[27],
		Sex2:         record[28],
		Treatment2:   record[29],
		Disease2:     record[30],
		Subject3:     record[31],
		Part3:        record[32],
		Sex3:         record[33],
		Treatment3:   record[34],
		Disease3:     record[35],
		Kana:         record[36],
		Condition:    record[37],
	}, nil
}
