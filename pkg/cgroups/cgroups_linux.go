package cgroups

import (
	"errors"
	"os"
	"strconv"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/taku-k/cgroups"
	"github.com/taku-k/log2s3-go/pkg"
)

type CgroupMngr struct {
	ctrl cgroups.Cgroup
	cfg  *pkg.UploadConfig
}

func NewCgroupMngr(c *pkg.UploadConfig) (*CgroupMngr, error) {
	if !isEnableLimit(c) {
		return nil, errors.New("No limit resources")
	}
	cpu := createCPULimit(c)
	memory := createMemoryLimit(c)
	network, minor := createNetCls(c)
	ctrl, err := cgroups.New(cgroups.V2, cgroups.StaticPath("/log2s3"), &specs.LinuxResources{
		CPU:     cpu,
		Memory:  memory,
		Network: network,
	})
	if err != nil {
		return nil, err
	}
	pid := os.Getpid()
	if err := ctrl.Add(cgroups.Process{Pid: pid}); err != nil {
		return nil, err
	}
	if isEnableLimitBW(c) {
		if err := createLimitBW(c, minor); err != nil {
			return nil, err
		}
	}
	return &CgroupMngr{ctrl, c}, nil
}

func createCPULimit(c *pkg.UploadConfig) *specs.LinuxCPU {
	var cpu *specs.LinuxCPU
	limit := c.CPULimit
	if limit != 0 {
		quota := int64(limit * 1000)
		cpu = &specs.LinuxCPU{
			Quota: &quota,
		}
	}
	return cpu
}

func createMemoryLimit(c *pkg.UploadConfig) *specs.LinuxMemory {
	var memory *specs.LinuxMemory
	limit := c.MemoryLimit
	if limit != 0 {
		memoryLimit := int64(limit * 1000 * 1000)
		memory = &specs.LinuxMemory{
			Limit: &memoryLimit,
		}
	}
	return memory
}

func createNetCls(c *pkg.UploadConfig) (*specs.LinuxNetwork, int) {
	var network *specs.LinuxNetwork
	minor := 1
	if isEnableLimitBW(c) {
		i32, _ := strconv.ParseInt("0x00100001", 16, 32)
		cls := uint32(i32)
		network = &specs.LinuxNetwork{
			ClassID: &cls,
		}
	}
	return network, minor
}

func isEnableLimit(c *pkg.UploadConfig) bool {
	return c.CPULimit > 0 || c.MemoryLimit > 0 || c.RateLimit > 0
}

func (c *CgroupMngr) Close() {
	if c.ctrl != nil {
		c.ctrl.Delete()
	}
	if isEnableLimitBW(c.cfg) {
		deleteLimitBW(c.cfg)
	}
}
