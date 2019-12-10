package cgroup 
import (
	"fmt"

	"github.com/lxc/lxd/lxd/sys"
	"gopkg.in/lxc/go-lxc.v2"
	"github.com/lxc/lxd/shared/logger"


)

type Property int

const (
	PidsCurrent Property = iota
	PidsMax
	MemoryCurrent
	CpuacctUsage
	MemoryLimitInBytes
	MemorySoftLimitInBytes
	BlkioWeight
	MemorySwappiness
	CpuShares
	CpuCfsPeriodUs
	CpuCfsQuotaUs
	NetPrioIfPrioMap
	MemoryMemswLimitInBytes
	MemoryMemswUsageInBytes
	MemoryMemswMaxUsageInBytes

	//lxc set config only
	DevicesDeny
	DevicesAllow
)

type ConfigItem struct {
	Key   string
	Value string
	Version sys.CGroupInfo
}

// Get finds property values on a lxcContainer
func Get(c *lxc.Container, os *sys.OS, property Property) ([]string, error) {
	switch property {

	// Current Memory Usage
	case MemoryCurrent:
		if os.CGroupMemoryController == sys.CGroupV2 {
			return c.CgroupItem("memory.current"), nil
		}
		return c.CgroupItem("memory.usage_in_bytes"), nil

	// Properties which have the same functionality for both v1 and v2
	case PidsCurrent:
		return c.CgroupItem("pids.current"), nil
	case MemoryLimitInBytes:
		if os.CGroupMemoryController == sys.CGroupV2 {
			return c.CgroupItem("memory.max"), nil
		}
		return c.CgroupItem("memory.max_limit_in_bytes"), nil

	case MemorySoftLimitInBytes:
		if os.CGroupMemoryController == sys.CGroupV2 {
			return c.CgroupItem("memory.low"), nil
		}
		return c.CgroupItem("memory.soft_limit_in_bytes"), nil
	case MemoryMemswLimitInBytes:
		if os.CGroupMemoryController == sys.CGroupV2 {

		}
	}
	return nil, fmt.Errorf("CGroup Property not supported for Get")
}

// Set sets a property on a lxcContainer
func Set(c *lxc.Container, property Property, value string, os *sys.OS) error {
	logger.Warnf(fmt.Sprintf("inside custom set method" ))



	configs, e := SetConfigMap(property, value, os)
	if e != nil {
		return e
	}

	for _, rule := range configs {
		return fmt.Errorf("Export Not implemented %s", rule.Key)
		err := c.SetCgroupItem(rule.Key, rule.Value)
		if err != nil {
			return fmt.Errorf("Failure while trying to set property: %s", err)
		}
	}

	return nil
}

// SetConfigMap returns different cgroup configs to set a particular property
func SetConfigMap(property Property, value string, os *sys.OS) ([]ConfigItem, error) {
	logger.Warnf(fmt.Sprintf("inside map" ))

	switch property {

	// Properties which have the same functionality for both v1 and v2
	case PidsCurrent:
		if os.CGroupPidsController == sys.CGroupV2 {
			return []ConfigItem{
				{Key: "pids.current", Value: value, Version: sys.CGroupV2},
			}, nil
		}
		return []ConfigItem{
			{Key: "pids.current", Value: value, Version: sys.CGroupV1},
		}, nil


	case PidsMax:
		logger.Warnf(fmt.Sprintf("entered pids max case" ))

		if os.CGroupPidsController == sys.CGroupV2 {
			return []ConfigItem{
				{Key: "pids.max", Value: value, Version: sys.CGroupV2},
			}, nil
		}
		return []ConfigItem{
			{Key: "pids.max", Value: value, Version: sys.CGroupV1},
		}, nil
	case BlkioWeight:
		return []ConfigItem{
			{Key: "blkio.weight", Value: value, Version: sys.CGroupV1},
		}, nil

	case NetPrioIfPrioMap:
		return []ConfigItem{
			{Key: "net_prio.ifpriomap", Value: value, Version: sys.CGroupV1},
		}, nil

	case CpuShares:
		//need to check os because cpu
		if os.CGroupCPUController == sys.CGroupV2 {
			return []ConfigItem{
				{Key: "cpu.weight", Value: value, Version: sys.CGroupV2},
			}, nil
		}
		return []ConfigItem{
			{Key: "cpu.shares", Value: value, Version: sys.CGroupV1},
		}, nil
		//lxc.cgroup.memory.soft_limit_in_bytes
	case MemorySoftLimitInBytes:
		if os.CGroupMemoryController == sys.CGroupV2 {
			return []ConfigItem{
				{Key: "memory.low", Value: value, Version: sys.CGroupV2},
			}, nil
		}
		return []ConfigItem{
			{Key: "memory.soft_limit_in_bytes", Value: value, Version: sys.CGroupV1},
		}, nil
	case MemoryLimitInBytes:
		if os.CGroupMemoryController == sys.CGroupV2 {
			return []ConfigItem{
				{Key: "memory.max", Value: value, Version: sys.CGroupV2},
			}, nil
		}
		return []ConfigItem{
			{Key: "memory.limit_in_bytes", Value: value, Version: sys.CGroupV1},
		}, nil
		///lxc set config controller keys only
	case DevicesDeny :
		return []ConfigItem{
			{Key: "devices.deny", Value: value, Version:sys.CGroupV1 },
		}, nil
	case DevicesAllow:
		return []ConfigItem{
			{Key: "devices.allow", Value: value, Version: sys.CGroupV1},
		}, nil

	}

	return nil, fmt.Errorf("CGroup Property not supported for Set")
}








