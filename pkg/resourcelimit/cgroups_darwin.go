package resourcelimit

import "github.com/hatena/u2s3/pkg/config"

type DummyCgroupMngr struct{}

func (c *DummyCgroupMngr) Close() {}

func NewCgroupMngr(c *config.UploadConfig) (*DummyCgroupMngr, error) {
	return &DummyCgroupMngr{}, nil
}
