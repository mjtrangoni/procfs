// Copyright 2017 The Prometheus Authors
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

package procfs

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// A ZoneInfo is the details parsed from /proc/zoneinfo (since Linux 2.6.13).
// > This file display information about memory zones. This is useful for
// > analyzing virtual memory behavior.
type ZoneInfo struct {
	Node string
	Zone string
	//Values []map[string]float64
}

// NewZoneInfo reads the zoneinfo statistics.
func NewZoneInfo() ([]ZoneInfo, error) {
	fs, err := NewFS(DefaultMountPoint)
	if err != nil {
		return nil, err
	}

	return fs.NewZoneInfo()
}

// NewZoneInfo reads the zoneinfo statistics from the specified `proc` filesystem.
func (fs FS) NewZoneInfo() ([]ZoneInfo, error) {
	file, err := os.Open(fs.Path("zoneinfo"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseZoneInfo(file)
}

func parseZoneInfo(r io.Reader) ([]ZoneInfo, error) {
	var (
		zoneInfo = []ZoneInfo{}
		scanner  = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Node") {
			parts := strings.Fields(line)

			node := strings.TrimRight(parts[1], ",")
			zone := strings.TrimRight(parts[3], ",")

			zoneInfo = append(zoneInfo, ZoneInfo{node, zone})
		}
	}

	return zoneInfo, scanner.Err()
}
