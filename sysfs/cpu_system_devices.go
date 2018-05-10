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

	"github.com/prometheus/common/log"
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
	KernelMax int64   // /sys/devices/system/cpu/kernel_max
	Offline   []int64 // /sys/devices/system/cpu/offline
	Online    []int64 // /sys/devices/system/cpu/online
	Possible  []int64 // /sys/devices/system/cpu/possible
	Present   []int64 // /sys/devices/system/cpu/present
}

// CPUInfo contains all CPU information.
type CPUInfo struct {
	CPUInfoGeneric CPUInfoGeneric
	CPUFreqString  []CPUFreq
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
		return cpuInformation, err
	}
	// Get CPUInfoGeneric information
	cpuInformation.CPUFreqString, err = parseCPUFreq(fs,
		cpuInformation.CPUInfoGeneric.Online)
	if err != nil {
		return cpuInformation, err
	}
	return cpuInformation, err
}

func parseCPUFreq(fs FS, online []int64) ([]CPUFreq, error) {

	cpuFreqStruct := make([]CPUFreq, len(online))
	var err error

	for _, cpuNum := range online {
		path := fs.Path("devices/system/cpu/cpu" + fmt.Sprintf("%d", cpuNum) + "/cpufreq")
		files, err := ioutil.ReadDir(path)
		if err != nil {
			// There is cases where there is no cpufreq information.
			continue
		}

		for _, fileDir := range files {

			fileContents, err := ioutil.ReadFile(path + "/" + fileDir.Name())
			if err != nil {
				return cpuFreqStruct, fmt.Errorf("cannot access %s, %s", path+"/"+fileDir.Name(), err)
			}
			value := strings.TrimSpace(string(fileContents))

			switch label := fileDir.Name(); label {
			case "cpuinfo_cur_freq":
				cpuFreqStruct[cpuNum].CPUInfoCurFreq, err = strconv.ParseInt(value, 10, 64)
			case "cpuinfo_max_freq":
				cpuFreqStruct[cpuNum].CPUInfoMaxFreq, err = strconv.ParseInt(value, 10, 64)
			case "cpuinfo_min_freq":
				cpuFreqStruct[cpuNum].CPUInfoMinFreq, err = strconv.ParseInt(value, 10, 64)
			case "cpuinfo_transition_latency":
				cpuFreqStruct[cpuNum].CPUInfoTransitionLatency, err = strconv.ParseInt(value, 10, 64)
			case "scaling_available_governors":
				cpuFreqStruct[cpuNum].ScalingAvailableGovernors = value
			case "scaling_cur_freq":
				cpuFreqStruct[cpuNum].ScalingCurFreq, err = strconv.ParseInt(value, 10, 64)
			case "scaling_driver":
				cpuFreqStruct[cpuNum].ScalingDriver = value
			case "scaling_governor":
				cpuFreqStruct[cpuNum].ScalingGovernor = value
			case "scaling_max_freq":
				cpuFreqStruct[cpuNum].ScalingMaxFreq, err = strconv.ParseInt(value, 10, 64)
			case "scaling_min_freq":
				cpuFreqStruct[cpuNum].ScalingMinFreq, err = strconv.ParseInt(value, 10, 64)
			case "scaling_setspeed":
				cpuFreqStruct[cpuNum].ScalingSetspeed, err = strconv.ParseInt(value, 10, 64)
			}
			if err != nil {
				log.Debugln(err)
			}
		}
	}
	return cpuFreqStruct, err
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
		fileContents, err := ioutil.ReadFile(path + "/" + fileDir.Name())
		if err != nil {
			return cpuInfoGeneric, fmt.Errorf("cannot access %s, %s", path+"/"+fileDir.Name(), err)
		}
		value := strings.TrimSpace(string(fileContents))

		switch label := fileDir.Name(); label {
		case "kernel_max":
			cpuInfoGeneric.KernelMax, err = strconv.ParseInt(value, 10, 64)
		case "offline":
			cpuInfoGeneric.Offline = parseCPURange(value)
		case "online":
			cpuInfoGeneric.Online = parseCPURange(value)
		case "possible":
			cpuInfoGeneric.Possible = parseCPURange(value)
		case "present":
			cpuInfoGeneric.Present = parseCPURange(value)
		}
	}
	return cpuInfoGeneric, err
}

func parseCPURange(value string) []int64 {
	var cpuSlice []int64
	for _, component := range strings.Split(value, ",") {
		first, err := strconv.ParseInt(strings.Split(component, "-")[0], 10, 64)
		if err != nil {
			log.Debugln(err)
		}
		last, err := strconv.ParseInt(strings.Split(component, "-")[1], 10, 64)
		if err != nil {
			log.Debugln(err)
		}
		for i := first; i <= last; i++ {
			cpuSlice = append(cpuSlice, i)
		}
	}
	return cpuSlice
}
