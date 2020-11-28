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
	if err := u.repo.Save(ctx, *r); err != nil {
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

	receipt := &model.Receipt{}
	var receiptNo uint32

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, xerrors.Errorf("on read.Read receipt file empty: %w", err)
		}

		switch record[0] {
		case model.IRRecordType:
			ir, err := ir(record)
			if err != nil {
				return nil, xerrors.Errorf("on read.ir couldn't read ir record: %w", err)
			}
			receipt.IR = *ir
		case model.RERecordType:
			re, err := re(record)
			if err != nil {
				return nil, xerrors.Errorf("on read.re couldn't read re record: %w", err)
			}
			receiptNo = re.ReceiptNo
			receipt.ReceiptItem(receiptNo).RE = *re
		case model.SYRecordType:
			sy, err := sy(record)
			if err != nil {
				return nil, xerrors.Errorf("on read.sy couldn't read sy record: %w", err)
			}
			receipt.ReceiptItem(receiptNo).SYs = append(receipt.ReceiptItem(receiptNo).SYs, *sy)
		case model.SIRecordType:
			si, err := si(record)
			if err != nil {
				return nil, xerrors.Errorf("on read.si couldn't read si record: %w", err)
			}
			receipt.ReceiptItem(receiptNo).SIInfos = append(receipt.ReceiptItem(receiptNo).SIInfos, model.SIInfo{SI: *si})
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

func re(record []string) (*model.RE, error) {
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

func sy(record []string) (*model.SY, error) {
	outcomeType, err := strconv.ParseUint(record[3], 10, 8)
	if err != nil {
		return nil, xerrors.Errorf("on sy.ParseUnit OutcomeType couldn't convert number from %v: %w", record[3], err)
	}
	return &model.SY{
		RecordType:  record[0],
		DiseaseID:   record[1],
		ReceiptedAt: record[2],
		OutcomeType: uint8(outcomeType),
		ModifierID:  record[4],
		DiseaseName: record[5],
		MainDisease: record[6],
		Comment:     record[7],
	}, nil
}

func si(record []string) (*model.SI, error) {
	times, err := strconv.ParseUint(record[6], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on si.ParseUnit Times couldn't convert number from %v: %w", record[6], err)
	}
	return &model.SI{
		RecordType:    record[0],
		TreatmentType: record[1],
		ChargeType:    record[2],
		TreatmentID:   record[3],
		Quantity:      record[4],
		Point:         record[5],
		Times:         uint16(times),
		CommentID1:    record[7],
		Comment1:      record[8],
		CommentID2:    record[9],
		Comment2:      record[10],
		CommentID3:    record[11],
		Comment3:      record[12],
		Day1:          record[13],
		Day2:          record[14],
		Day3:          record[15],
		Day4:          record[16],
		Day5:          record[17],
		Day6:          record[18],
		Day7:          record[19],
		Day8:          record[20],
		Day9:          record[21],
		Day10:         record[22],
		Day11:         record[23],
		Day12:         record[24],
		Day13:         record[25],
		Day14:         record[26],
		Day15:         record[27],
		Day16:         record[28],
		Day17:         record[29],
		Day18:         record[30],
		Day19:         record[31],
		Day20:         record[32],
		Day21:         record[33],
		Day22:         record[34],
		Day23:         record[35],
		Day24:         record[36],
		Day25:         record[37],
		Day26:         record[38],
		Day27:         record[39],
		Day28:         record[40],
		Day29:         record[41],
		Day30:         record[42],
		Day31:         record[43],
	}, nil
}
