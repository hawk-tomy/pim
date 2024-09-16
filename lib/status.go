/*
MIT License

# Copyright (c) 2023 - present hawk-tomy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package lib

import (
	"fmt"
	"github.com/hawk-tomy/pim/lib/list"
	"sort"
)

var (
	installedPythonVersions map[int]Version
	updatablePythonVersions map[int]*list.Element[Version]
)

func getInstalledPythonVersions(config Config) error {
	installedPythonVersions = make(map[int]Version)
	for _, v := range getPythonVersions(config) {
		installedPythonVersions[v.Minor] = v
	}
	return nil
}

func detectUpdatablePythonVersions(config Config) error {
	if err := fetchLatestVersions(config); err != nil {
		return err
	}

	if err := getInstalledPythonVersions(config); err != nil {
		return err
	}

	updatablePythonVersions = make(map[int]*list.Element[Version])
	for _, installedVersion := range installedPythonVersions {
		// same minor version and installed version is old version
		if v, ok := fetchedVersions[installedVersion.Minor]; ok {
			ver := v.Back()
			for ver != nil {
				if ver.Value.LessThanOrEqual(installedVersion) {
					ver = nil
					break
				}
				if v_, ok_ := failedMinimumVersions[installedVersion.Minor]; ok_ && v_.LessThanOrEqual(ver.Value) {
					ver = ver.Prev()
					continue // can not update.
				}
				break
			}
			if ver == nil {
				continue
			}
			updatablePythonVersions[installedVersion.Minor] = ver
		}
	}

	return nil
}

func StatusCommand(config Config) error {
	if err := detectUpdatablePythonVersions(config); err != nil {
		return err
	}

	MinorVersions := make([]int, 0, len(installedPythonVersions))
	for k := range installedPythonVersions {
		MinorVersions = append(MinorVersions, k)
	}
	sort.Ints(MinorVersions)

	fmt.Println("Installed Python versions:")
	for _, minor := range MinorVersions {
		v := installedPythonVersions[minor]
		statusStr := v.String()
		if ver, ok := updatablePythonVersions[minor]; ok {
			statusStr += fmt.Sprintf(" (updatable: %s)", ver.Value.String())
		}
		fmt.Println(statusStr)
	}

	return nil
}
