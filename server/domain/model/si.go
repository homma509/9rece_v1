package model

const (
	// SIRecordType 診療行為レコードのレコード識別情報
	SIRecordType = "SI"
)

// SI 診療行為レコード
type SI struct {
	RecordType    string // レコード識別情報
	TreatmentType string // 診療識別
	ChargeType    string // 負担区分
	TreatmentID   string // 診療行為コード
	Quantity      string // 数量データ
	Point         string // 点数
	Times         uint16 // 回数
	CommentID1    string // コメントコード
	Comment1      string // 文字データ
	CommentID2    string // コメントコード
	Comment2      string // 文字データ
	CommentID3    string // コメントコード
	Comment3      string // 文字データ
	Day1          string // 1日の情報
	Day2          string // 2日の情報
	Day3          string // 3日の情報
	Day4          string // 4日の情報
	Day5          string // 5日の情報
	Day6          string // 6日の情報
	Day7          string // 7日の情報
	Day8          string // 8日の情報
	Day9          string // 9日の情報
	Day10         string // 10日の情報
	Day11         string // 11日の情報
	Day12         string // 12日の情報
	Day13         string // 13日の情報
	Day14         string // 14日の情報
	Day15         string // 15日の情報
	Day16         string // 16日の情報
	Day17         string // 17日の情報
	Day18         string // 18日の情報
	Day19         string // 19日の情報
	Day20         string // 20日の情報
	Day21         string // 21日の情報
	Day22         string // 22日の情報
	Day23         string // 23日の情報
	Day24         string // 24日の情報
	Day25         string // 25日の情報
	Day26         string // 26日の情報
	Day27         string // 27日の情報
	Day28         string // 28日の情報
	Day29         string // 29日の情報
	Day30         string // 30日の情報
	Day31         string // 31日の情報
}
