package resourcelimit

import (
	"errors"
	"os"
	"strconv"

	"github.com/containerd/cgroups"
	"github.com/hatena/u2s3/pkg/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type CgroupMngr struct {
	ctrl cgroups.Cgroup
	cfg  *config.UploadConfig
}

func NewCgroupMngr(c *config.UploadConfig) (*CgroupMngr, error) {
	if !isEnableLimit(c) {
		return nil, errors.New("No limit resources")
	}
	cpu := createCPULimit(c)
	memory := createMemoryLimit(c)
	network, minor := createNetCls(c)
	ctrl, err := cgroups.New(cgroups.V1, cgroups.StaticPath("/u2s3"), &specs.LinuxResources{
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

func createCPULimit(c *config.UploadConfig) *specs.LinuxCPU {
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

func createMemoryLimit(c *config.UploadConfig) *specs.LinuxMemory {
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

func createNetCls(c *config.UploadConfig) (*specs.LinuxNetwork, int) {
	var network *specs.LinuxNetwork
	minor := 1
	if isEnableLimitBW(c) {
		i32, _ := strconv.ParseInt("00100001", 16, 32)
		cls := uint32(i32)
		network = &specs.LinuxNetwork{
			ClassID: &cls,
		}
	}
	return network, minor
}

func isEnableLimit(c *config.UploadConfig) bool {
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
