package resourcelimit

import "github.com/taku-k/u2s3/pkg"

type DummyCgroupMngr struct{}

func (c *DummyCgroupMngr) Close() {}

func NewCgroupMngr(c *pkg.UploadConfig) (*DummyCgroupMngr, error) {
	return &DummyCgroupMngr{}, nil
}
