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
	"sort"
)

func UpdateLatest(config Config, version Version) error {
	if err := detectUpdatablePythonVersions(config); err != nil {
		return err
	}
	minor := version.Minor
	if latestVersion, ok := updatablePythonVersions[minor]; ok {
		return doUpdate(config, latestVersion)
	}
	return fmt.Errorf("this version is already latest version")
}

func UpdateAll(config Config) error {
	if err := detectUpdatablePythonVersions(config); err != nil {
		return err
	}

	if len(updatablePythonVersions) == 0 {
		fmt.Println("There is no updatable python.")
		return nil
	}

	minorVersions := make([]int, 0, len(updatablePythonVersions))
	for k := range updatablePythonVersions {
		minorVersions = append(minorVersions, k)
	}
	sort.Ints(minorVersions)

	fmt.Printf("update versions are:\n")
	for _, minor := range minorVersions {
		var verStr string
		if v, ok := installedPythonVersions[minor]; ok {
			verStr = v.String()
		}
		fmt.Printf("%s -> %s\n", verStr, updatablePythonVersions[minor].Value.String())
	}

	if !Confirm("Do you want to update all updatable python? [Y/n]") {
		return nil
	}

	fmt.Println("start updating...")
	for _, minor := range minorVersions {
		ver := updatablePythonVersions[minor]
		fmt.Printf("updating python %s\n", ver.Value.String())
		if err := doUpdate(config, ver); err != nil {
			fmt.Printf("an error occurred while updating python %s: %s\n", ver.Value.String(), err.Error())
		}
	}

	return nil
}
