package aggregator

import (
	"github.com/taku-k/log2s3-go/pkg/aws"
	"github.com/taku-k/log2s3-go/pkg"
)

type Uploader struct {
	cli *aws.S3Cli
}

func NewUploader(config *pkg.UploadConfig) *Uploader {
	cli := aws.NewS3ForTest(config)
	return &Uploader{cli}
}

func (u *Uploader) Upload(e *Epoch) error {
	e.writer.Close()
	e.buf.Flush()

	seq := 0
	key := ""
	var err error
	for {
		seq += 1
		key, err = e.GetObjectKey(seq)
		if err != nil {
			return err
		}
		if !u.cli.HasKey(key) {
			break
		}
	}
	if err := u.cli.Upload(key, e.fp); err != nil {
		return err
	}
	return nil
}
