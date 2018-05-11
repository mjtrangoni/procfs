// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sysfs

import (
	"testing"
)

func TestNewCPUInfo(t *testing.T) {
	fs, err := NewFS("fixtures")
	if err != nil {
		t.Fatal(err)
	}

	nc, err := fs.NewCPUInfo()
	if err != nil {
		t.Fatal(err)
	}

	// Test CPUInfoGeneric
	if want, got := int64(2047), nc.CPUInfoGeneric.KernelMax; want != got {
		t.Errorf("want kernel_max %d, got %d", want, got)
	}
	if want, got := int64(89), nc.CPUInfoGeneric.Online[23]; want != got {
		t.Errorf("want online %d, got %d", want, got)
	}

	// Test CPUFreq cpu0 - AMD Epic CentOS7.4.1708 x86_64
	if want, got := int64(2300000), nc.CPUFreqSlice[0].CPUInfoCurFreq; want != got {
		t.Errorf("want cpu0/cpufreq/cpuinfo_cur_freq %d, got %d", want, got)
	}
	if want, got := "acpi-cpufreq", nc.CPUFreqSlice[0].ScalingDriver; want != got {
		t.Errorf("want cpu0/cpufreq/scaling_driver %s, got %s", want, got)
	}
	if want, got := int64(0), nc.CPUFreqSlice[0].ScalingSetspeed; want != got {
		t.Errorf("want cpu0/cpufreq/scaling_setspeed %d, got %d", want, got)
	}

	// Test CPUFreq cpu1 - Intel Skylake CentOS7.4.1708 x86_64
	if want, got := int64(1178125), nc.CPUFreqSlice[1].CPUInfoCurFreq; want != got {
		t.Errorf("want cpu1/cpufreq/cpuinfo_cur_freq %d, got %d", want, got)
	}
	if want, got := "intel_pstate", nc.CPUFreqSlice[1].ScalingDriver; want != got {
		t.Errorf("want cpu1/cpufreq/scaling_driver %s, got %s", want, got)
	}
	if want, got := int64(0), nc.CPUFreqSlice[1].ScalingSetspeed; want != got {
		t.Errorf("want cpu1/cpufreq/scaling_setspeed %d, got %d", want, got)
	}

	// Test CPUFreq cpu8 - IBM Power8 CentOS7.4.1708 ppc64le
	if want, got := int64(3690000), nc.CPUFreqSlice[8].CPUInfoCurFreq; want != got {
		t.Errorf("want cpu8/cpufreq/cpuinfo_cur_freq %d, got %d", want, got)
	}
	if want, got := "powernv-cpufreq", nc.CPUFreqSlice[8].ScalingDriver; want != got {
		t.Errorf("want cpu8/cpufreq/scaling_driver %s, got %s", want, got)
	}
	if want, got := int64(3690000), nc.CPUFreqSlice[8].ScalingCurFreq; want != got {
		t.Errorf("want cpu8/cpufreq/scaling_cur_freq %d, got %d", want, got)
	}
	if want, got := int64(0), nc.CPUFreqSlice[8].ScalingSetspeed; want != got {
		t.Errorf("want cpu8/cpufreq/scaling_setspeed %d, got %d", want, got)
	}

	// Test CPUTopology cpu0 - AMD Epic CentOS7.4.1708 x86_64
	if want, got := int64(0), nc.CPUTopologySlice[0].CoreID; want != got {
		t.Errorf("want cpu0/topology/core_id %d, got %d", want, got)
	}
	if want, got := "0-23,48-71", nc.CPUTopologySlice[0].CoreSiblingsList; want != got {
		t.Errorf("want cpu0/topology/core_siblings_list %s, got %s", want, got)
	}
	// Test CPUTopology cpu1 - Intel Skylake CentOS7.4.1708 x86_64
	if want, got := int64(1), nc.CPUTopologySlice[1].CoreID; want != got {
		t.Errorf("want cpu1/topology/core_id %d, got %d", want, got)
	}
	if want, got := "0-15,32-47", nc.CPUTopologySlice[1].CoreSiblingsList; want != got {
		t.Errorf("want cpu1/topology/core_siblings_list %s, got %s", want, got)
	}

	// Test CPUTopology cpu8 - IBM Power8 CentOS7.4.1708 ppc64le
	if want, got := int64(48), nc.CPUTopologySlice[8].CoreID; want != got {
		t.Errorf("want cpu8/topology/core_id %d, got %d", want, got)
	}
	if want, got := "0-1,8-9,16-17,24-25,32-33", nc.CPUTopologySlice[8].CoreSiblingsList; want != got {
		t.Errorf("want cpu8/topology/core_siblings_list %s, got %s", want, got)
	}

	// Test CPUThermalThrottle cpu0 - AMD Epic CentOS7.4.1708 x86_64 - nothing there
	if want, got := int64(0), nc.CPUThermalThrottleSlice[0].CoreThrottleCount; want != got {
		t.Errorf("want cpu0/thermal_throttle/core_throttle_count %d, got %d", want, got)
	}
	if want, got := int64(0), nc.CPUThermalThrottleSlice[0].PackageThrottleCount; want != got {
		t.Errorf("want cpu0/thermal_throttle/package_throttle_count %d, got %d", want, got)
	}
	// Test CPUThermalThrottle cpu1 - Intel Skylake CentOS7.4.1708 x86_64
	if want, got := int64(30), nc.CPUThermalThrottleSlice[1].CoreThrottleCount; want != got {
		t.Errorf("want cpu0/thermal_throttle/core_throttle_count %d, got %d", want, got)
	}
	if want, got := int64(45), nc.CPUThermalThrottleSlice[1].PackageThrottleCount; want != got {
		t.Errorf("want cpu0/thermal_throttle/package_throttle_count %d, got %d", want, got)
	}
}
