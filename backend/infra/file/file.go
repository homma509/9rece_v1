package file

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
			return err
		}
		f.client = s3.New(sess)
	}
	return nil
}

// GetObject ファイルを取得します
func (f *File) GetObject(bucket, key string) (io.ReadCloser, error) {
	err := f.connect()
	if err != nil {
		return nil, err
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

// MoveObject ファイルを移動します
func (f *File) MoveObject(bucket, src, dst string) error {
	err := f.connect()
	if err != nil {
		return err
	}

	_, err = f.client.CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(fmt.Sprintf("%s/%s", bucket, src)),
		Key:        aws.String(dst),
	})
	if err != nil {
		return err
	}

	_, err = f.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(src),
	})
	if err != nil {
		return err
	}

	return nil
}
