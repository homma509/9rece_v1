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

	record, err := r.Read()
	if err == io.EOF {
		return nil, xerrors.Errorf("on read.Read receipt file EOF: %w", err)
	}
	if err != nil {
		return nil, xerrors.Errorf("on read.Read receipt file empty: %w", err)
	}

	receipt := &model.Receipt{
		ReceiptItems: map[uint32]*model.ReceiptItem{},
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
	benefitRate, err := strconv.ParseUint(record[7], 10, 8)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit BenefitRate couldn't convert number from %v: %w", record[7], err)
	}
	admittedAd, err := strconv.ParseUint(record[8], 10, 32)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit AdmittedAd couldn't convert number from %v: %w", record[8], err)
	}
	chargeType, err := strconv.ParseUint(record[10], 10, 8)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit ChargeType couldn't convert number from %v: %w", record[10], err)
	}
	wardCount, err := strconv.ParseUint(record[12], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit WardCount couldn't convert number from %v: %w", record[12], err)
	}
	discountUnit, err := strconv.ParseUint(record[14], 10, 8)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit DiscountUnit couldn't convert number from %v: %w", record[14], err)
	}
	reserve1, err := strconv.ParseUint(record[15], 10, 8)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Reserve1 couldn't convert number from %v: %w", record[15], err)
	}
	reserve2, err := strconv.ParseUint(record[16], 10, 8)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Reserve2 couldn't convert number from %v: %w", record[16], err)
	}
	reserve3, err := strconv.ParseUint(record[17], 10, 8)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Reserve3 couldn't convert number from %v: %w", record[17], err)
	}
	reserve4, err := strconv.ParseUint(record[19], 10, 32)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Reserve4 couldn't convert number from %v: %w", record[19], err)
	}
	part1, err := strconv.ParseUint(record[22], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Part1 couldn't convert number from %v: %w", record[22], err)
	}
	sex1, err := strconv.ParseUint(record[23], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Sex1 couldn't convert number from %v: %w", record[23], err)
	}
	treatment1, err := strconv.ParseUint(record[24], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Treatment1 couldn't convert number from %v: %w", record[24], err)
	}
	disease1, err := strconv.ParseUint(record[25], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Disease1 couldn't convert number from %v: %w", record[25], err)
	}
	part2, err := strconv.ParseUint(record[27], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Part2 couldn't convert number from %v: %w", record[27], err)
	}
	sex2, err := strconv.ParseUint(record[28], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Sex2 couldn't convert number from %v: %w", record[28], err)
	}
	treatment2, err := strconv.ParseUint(record[29], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Treatment2 couldn't convert number from %v: %w", record[29], err)
	}
	disease2, err := strconv.ParseUint(record[30], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Disease2 couldn't convert number from %v: %w", record[30], err)
	}
	part3, err := strconv.ParseUint(record[32], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Part3 couldn't convert number from %v: %w", record[32], err)
	}
	sex3, err := strconv.ParseUint(record[33], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Sex3 couldn't convert number from %v: %w", record[33], err)
	}
	treatment3, err := strconv.ParseUint(record[34], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Treatment3 couldn't convert number from %v: %w", record[34], err)
	}
	disease3, err := strconv.ParseUint(record[35], 10, 16)
	if err != nil {
		return nil, xerrors.Errorf("on ir.ParseUnit Disease3 couldn't convert number from %v: %w", record[35], err)
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
		BenefitRate:  uint16(benefitRate),
		AdmittedAd:   uint32(admittedAd),
		WardType:     record[9],
		ChargeType:   uint8(chargeType),
		ReceiptNote:  record[11],
		WardCount:    uint16(wardCount),
		KarteNo:      record[13],
		DiscountUnit: uint8(discountUnit),
		Reserve1:     uint8(reserve1),
		Reserve2:     uint8(reserve2),
		Reserve3:     uint8(reserve3),
		SearchNo:     record[18],
		Reserve4:     uint32(reserve4),
		InvoiceInfo:  record[20],
		Subject1:     record[21],
		Part1:        uint16(part1),
		Sex1:         uint16(sex1),
		Treatment1:   uint16(treatment1),
		Disease1:     uint16(disease1),
		Subject2:     record[26],
		Part2:        uint16(part2),
		Sex2:         uint16(sex2),
		Treatment2:   uint16(treatment2),
		Disease2:     uint16(disease2),
		Subject3:     record[31],
		Part3:        uint16(part3),
		Sex3:         uint16(sex3),
		Treatment3:   uint16(treatment3),
		Disease3:     uint16(disease3),
		Kana:         record[36],
		Condition:    record[37],
	}, nil
}
