package file

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

// File ファイル構造体
type File struct {
	client *s3.S3
	config *aws.Config
}

// NewFile ファイルインターフェースを生成します
func NewFile(config *aws.Config) *File {
	return &File{
		config: config,
	}
}

func (f *File) connect() error {
	if f.client == nil {
		sess, err := session.NewSession(f.config)
		if err != nil {
			return errors.WithStack(err)
		}
		f.client = s3.New(sess)
	}
	return nil
}

// GetObject ファイルを取得します
func (f *File) GetObject(bucket, key string) (io.ReadCloser, error) {
	err := f.connect()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	obj, err := f.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return obj.Body, nil
}
