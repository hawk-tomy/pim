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

import "fmt"

func getInstalledPythonByMinor(config Config, version Version) (Version, error) {
	err := getInstalledPythonVersions(config)
	if err != nil {
		return Version{}, err
	}
	for _, v := range installedPythonVersions {
		if v.Major == version.Major && v.Minor == version.Minor {
			return v, nil
		}
	}
	return Version{}, fmt.Errorf("not found installed python: %s", version.String())
}

func UninstallPython(config Config, version Version) error {
	installedVersion, err := getInstalledPythonByMinor(config, version)
	if !Confirm(fmt.Sprintf("uninstall python %s? [Y/n]", installedVersion.String())) {
		return nil
	}

	fmt.Printf("uninstalling python %s\n", installedVersion.String())
	if err != nil {
		return err
	}
	return doUninstall(config, installedVersion)
}
