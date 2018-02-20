package core

import (
	"log"
	"os"

	"github.com/hatena/u2s3/pkg/aws"
	"github.com/hatena/u2s3/pkg/config"
)

type UploadableFile interface {
	GetObjectKey() (string, error)
	GetFile() *os.File
	ResetSeq()
	Flush()
}

type Uploader struct {
	cli *aws.S3Cli
}

func NewUploader(config *config.UploadConfig) *Uploader {
	cli := aws.NewS3Cli(config)
	return &Uploader{cli}
}

func (u *Uploader) Upload(uf UploadableFile) error {
	uf.Flush()

	key := ""
	var err error
	uf.ResetSeq()
	for {
		key, err = uf.GetObjectKey()
		if err != nil {
			return err
		}
		if !u.cli.HasKey(key) {
			break
		}
	}
	if err := u.cli.Upload(key, uf.GetFile()); err != nil {
		return err
	}
	log.Printf("[info] Uploaded %s\n", key)
	return nil
}
