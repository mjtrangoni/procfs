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
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// CPUFreq contains information about CPU frequency and voltage scaling. See
// https://www.kernel.org/doc/Documentation/cpu-freq/user-guide.txt
type CPUFreq struct {
	CPUInfoCurFreq            int64  // /sys/devices/system/cpu/cpu*/cpufreq/cpuinfo_cur_freq
	CPUInfoMaxFreq            int64  // /sys/devices/system/cpu/cpu*/cpufreq/cpuinfo_max_freq
	CPUInfoMinFreq            int64  // /sys/devices/system/cpu/cpu*/cpufreq/cpuinfo_min_freq
	CPUInfoTransitionLatency  int64  // /sys/devices/system/cpu/cpu*/cpufreq/cpuinfo_transition_latency
	ScalingAvailableGovernors string // /sys/devices/system/cpu/cpu*/cpufreq/scaling_available_governors
	ScalingCurFreq            int64  // /sys/devices/system/cpu/cpu*/cpufreq/scaling_cur_freq
	ScalingDriver             string // /sys/devices/system/cpu/cpu*/cpufreq/scaling_driver
	ScalingGovernor           string // /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor
	ScalingMaxFreq            int64  // /sys/devices/system/cpu/cpu*/cpufreq/scaling_max_freq
	ScalingMinFreq            int64  // /sys/devices/system/cpu/cpu*/cpufreq/scaling_min_freq
	ScalingSetspeed           int64  // /sys/devices/system/cpu/cpu*/cpufreq/scaling_setspeed
}

// CPUTopology contains information about the CPU topology. See
// https://www.kernel.org/doc/Documentation/cputopology.txt
type CPUTopology struct {
	CoreID             int64  // /sys/devices/system/cpu/cpu*/topology/core_id
	CoreSiblings       string // /sys/devices/system/cpu/cpu*/topology/core_siblings
	CoreSiblingsList   string // /sys/devices/system/cpu/cpu*/topology/core_siblings_list
	PhysicalPackageID  int64  // /sys/devices/system/cpu/cpu*/topology/physical_package_id
	ThreadSiblings     string // /sys/devices/system/cpu/cpu*/topology/thread_siblings
	ThreadSiblingsList string // /sys/devices/system/cpu/cpu*/topology/thread_siblings_list
}

// CPUThermalThrottle contains information about the CPU thermal throttling. See
// https://www.kernel.org/doc/Documentation/ ??
type CPUThermalThrottle struct {
	CoreThrottleCount    int64 // /sys/devices/system/cpu/cpu*/thermal_throttle/core_throttle_count
	PackageThrottleCount int64 // /sys/devices/system/cpu/cpu*/thermal_throttle/package_throttle_count
}

// CPUInfoGeneric contains information about all CPUs in general. See
// https://www.kernel.org/doc/Documentation/ABI/testing/sysfs-devices-system-cpu
type CPUInfoGeneric struct {
	KernelMax int64  // /sys/devices/system/cpu/kernel_max
	Offline   string // /sys/devices/system/cpu/offline
	Online    string // /sys/devices/system/cpu/online
	Possible  string // /sys/devices/system/cpu/possible
	Present   string // /sys/devices/system/cpu/present
}

// CPUInfo contains all CPU information.
type CPUInfo struct {
	CPUInfoGeneric CPUInfoGeneric
}

// NewCPUInfo reads the cpu information.
func NewCPUInfo() (CPUInfo, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return CPUInfo{}, err
	}

	return fs.NewCPUInfo()
}

// NewCPUInfo reads the cpu information from sysfs files.
func (fs FS) NewCPUInfo() (CPUInfo, error) {

	var err error
	cpuInformation := CPUInfo{}

	// Get CPUInfoGeneric information
	cpuInformation.CPUInfoGeneric, err = parseCPUInfoGeneric(fs)
	if err != nil {
		return CPUInfo{}, err
	}
	return cpuInformation, err
}

func parseCPUInfoGeneric(fs FS) (CPUInfoGeneric, error) {

	cpuInfoGeneric := CPUInfoGeneric{}

	path := fs.Path("devices/system/cpu")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return cpuInfoGeneric, fmt.Errorf("cannot access %s dir %s", path, err)
	}

	for _, fileDir := range files {
		// skip directories
		if fileDir.IsDir() {
			continue
		}
		fileContents, err := sysReadFile(path + "/" + fileDir.Name())
		if err != nil {
			return cpuInfoGeneric, fmt.Errorf("cannot access %s, %s", path+"/"+fileDir.Name(), err)
		}
		value := strings.TrimSpace(string(fileContents))

		switch label := fileDir.Name(); label {
		case "kernel_max":
			cpuInfoGeneric.KernelMax, err = strconv.ParseInt(value, 10, 64)
		case "offline":
			cpuInfoGeneric.Offline = value
		case "online":
			cpuInfoGeneric.Online = value
		case "possible":
			cpuInfoGeneric.Possible = value
		case "present":
			cpuInfoGeneric.Present = value
		}
	}
	return cpuInfoGeneric, err
}
