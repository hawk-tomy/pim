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
	"errors"
	"fmt"
)

func InstallPython(config Config, version Version, needLatest bool) error {
	if err := fetchLatestVersions(config); err != nil {
		return err
	}

	if needLatest {
		v := fetchedVersions[version.Minor].Back()
		for v != nil {
			if config.AllowPreRelease || v.Value.Pre == 0 {
				return doInstall(config, v, needLatest)
			}
			v = v.Prev()
		}
		return fmt.Errorf("not found latest version")
	}

	if versions, ok := fetchedVersions[version.Minor]; ok {
		v := versions.Front()
		for v != nil {
			if v.Value.Equal(version) {
				return doInstall(config, v, needLatest)
			}
			v = v.Next()
		}
	}
	return errors.New("the version is not found")
}
