/*
Copyright 2019 The MayaData Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import "github.com/pkg/errors"

// RaidGroupConfig contains raid type and device(s)
// count for a raid group
type RaidGroupConfig struct {
	// Type is the raid group type
	// Supported values are : stripe, mirror, raidz and raidz2
	RAIDType PoolRAIDType `json:"raidType"`
	// GroupDeviceCount contains device count in a raid group
	// -- for stripe DeviceCount = 1
	// -- for mirror DeviceCount = 2
	// -- for raidz DeviceCount = (2^n + 1) default is (2 + 1)
	// -- for raidz2 DeviceCount = (2^n + 2) default is (4 + 2)
	GroupDeviceCount int64 `json:"groupDeviceCount"`
}

// GetDefaultRaidGroupConfig returns an object of RaidGroupConfig
// with default configuration. Default raid group config is
// type: mirror groupDeviceCount: 2
func GetDefaultRaidGroupConfig() *RaidGroupConfig {
	return &RaidGroupConfig{
		RAIDType:         PoolRAIDTypeDefault,
		GroupDeviceCount: RAIDTypeToDefaultMinDiskCount[PoolRAIDTypeMirror],
	}
}

// PopulateDefaultGroupDeviceCountIfNotPresent populate default device count for
// a given raid group if device count is not set then. If device count for a raid
// group present then it will skip.
func (rgc *RaidGroupConfig) PopulateDefaultGroupDeviceCountIfNotPresent() error {
	if rgc.GroupDeviceCount != 0 {
		return nil
	}
	dc, ok := RAIDTypeToDefaultMinDiskCount[rgc.RAIDType]
	if !ok {
		return errors.Errorf("Invalid RAID type %q: Supports %q, %q, %q or %q.",
			rgc.RAIDType, PoolRAIDTypeStripe, PoolRAIDTypeMirror, PoolRAIDTypeRAIDZ, PoolRAIDTypeRAIDZ2)
	}
	rgc.GroupDeviceCount = dc
	return nil
}

// GetDataDeviceCount returns data device count for one raid group configuration
// This is helpfull to calculate pool capacity.
func (rgc *RaidGroupConfig) GetDataDeviceCount() int64 {
	switch rgc.RAIDType {
	case PoolRAIDTypeMirror:
		return rgc.GroupDeviceCount / 2
	// For stripe pool data device count in is n.
	// where n block devices present in raid group config
	case PoolRAIDTypeStripe:
		return rgc.GroupDeviceCount
	// For raidz pool data device count in is x - 1.
	// where x = 2^n + 1 block devices present in raid group config
	case PoolRAIDTypeRAIDZ:
		return rgc.GroupDeviceCount - 1
	// For raidz2 pool data device count in is x - 2.
	// where x = 2^n + 2 block devices present in raid group config
	case PoolRAIDTypeRAIDZ2:
		return rgc.GroupDeviceCount - 2
	default:
		return 0
	}
}

// Validate validates RaidGroupConfig
func (rgc *RaidGroupConfig) Validate() error {
	// If we got any -ve number or 0 then it an invalid device count.
	if rgc.GroupDeviceCount <= 0 {
		return errors.Errorf("Invalid device count %d for RAID type %q.",
			rgc.GroupDeviceCount, rgc.RAIDType)
	}

	minDeviceCount, ok := RAIDTypeToDefaultMinDiskCount[PoolRAIDType(rgc.RAIDType)]
	if !ok {
		return errors.Errorf("Invalid RAID type %q: Supports %q, %q, %q or %q.",
			rgc.RAIDType, PoolRAIDTypeStripe, PoolRAIDTypeMirror, PoolRAIDTypeRAIDZ, PoolRAIDTypeRAIDZ2)
	}

	// If device count is less than min device count then that is not a valid
	// raid group config. ie - For raid-z(2^n + 1) device count 1 is valid if
	// n=0 but that is not a valid raid group config. For raid-z2(2^n + 2) device
	// count 2 is valid if n=0 but that is not a valid raid group config.
	if minDeviceCount > rgc.GroupDeviceCount {
		return errors.Errorf("Invalid device count %d for RAID type %q",
			rgc.GroupDeviceCount, rgc.RAIDType)
	}

	switch rgc.RAIDType {
	// For mirror pool device count in one vdev is 2
	case PoolRAIDTypeMirror:
		{
			if rgc.GroupDeviceCount != 2 {
				return errors.Errorf("Invalid device count %d for RAID type %q: Want 2.",
					rgc.GroupDeviceCount, rgc.RAIDType)
			}
		}
	// For stripe pool device count in one vdev is n. Where n > 0
	case PoolRAIDTypeStripe:
		{
			return nil
		}
	// For raidz raid group device count in one vdev is (2^n +1)
	case PoolRAIDTypeRAIDZ:
		{
			count := rgc.GroupDeviceCount - 1
			for count != 1 {
				r := count % 2
				if r != 0 {
					return errors.Errorf("Invalid device count %d for RAID type %q: Want 2^n + 1.",
						rgc.GroupDeviceCount, rgc.RAIDType)
				}
				count = count / 2
			}
		}
	// For raidz2 raid group device count in one vdev is (2^n +1)
	case PoolRAIDTypeRAIDZ2:
		{
			count := rgc.GroupDeviceCount - 2
			for count != 1 {
				r := count % 2
				if r != 0 {
					return errors.Errorf("Invalid device count %d for RAID type %q: Want 2^n + 2.",
						rgc.GroupDeviceCount, rgc.RAIDType)
				}
				count = count / 2
			}
		}
	default:
		{
			return errors.Errorf("Invalid RAID type %q: Supports %q, %q, %q or %q.",
				rgc.RAIDType, PoolRAIDTypeStripe, PoolRAIDTypeMirror, PoolRAIDTypeRAIDZ, PoolRAIDTypeRAIDZ2)
		}
	}
	return nil
}
