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
	"regexp"
	"strconv"
	"strings"
)

const (
	regex = `^v?(?P<major>\d)+(?:\.(?P<minor>\d+))?(?:\.(?P<micro>\d+))?(?:(?P<pre>a|b|rc)(?P<preN>\d+))?$`
)

var (
	versionRegex *regexp.Regexp
)

func init() {
	versionRegex = regexp.MustCompile(regex)
}

type Version struct {
	Major  int
	Minor  int
	Micro  int
	Pre    int
	PreNum int
}

func NewVersion(versionString string) (Version, error) {
	var v Version
	var n int
	var err error

	matches := versionRegex.FindStringSubmatch(versionString)
	if matches == nil {
		return Version{}, fmt.Errorf("invalid version string (not match): %s", versionString)
	}

	l := len(matches)

	n = 0
	if l > 1 && matches[1] != "" {
		if n, err = strconv.Atoi(matches[1]); err != nil {
			return Version{}, fmt.Errorf("invalid version string (Major): %s", versionString)
		}
	}
	v.Major = n

	n = 0
	if l > 2 && matches[2] != "" {
		if n, err = strconv.Atoi(matches[2]); err != nil {
			return Version{}, fmt.Errorf("invalid version string (Minor): %s", versionString)
		}
	}
	v.Minor = n

	n = 0
	if l > 3 && matches[3] != "" {
		n, err = strconv.Atoi(matches[3])
		if err != nil {
			return Version{}, fmt.Errorf("invalid version string (Micro): %s", versionString)
		}
	}
	v.Micro = n

	if l > 4 && matches[4] != "" {
		if matches[4] == "rc" {
			v.Pre = -1
		} else if matches[4] == "b" {
			v.Pre = -2
		} else if matches[4] == "a" {
			v.Pre = -3
		} else {
			panic("UNREACHABLE: invalid version string (Pre)")
		}

		if n, err = strconv.Atoi(matches[5]); err != nil {
			return Version{}, fmt.Errorf("invalid version string (Pre): %s", versionString)
		}
		v.PreNum = n
	} else {
		v.Pre = 0
		v.PreNum = 0
	}
	return v, nil
}

func (v Version) Compare(o *Version) int {
	if v.Major > o.Major {
		return 1
	} else if v.Major < o.Major {
		return -1
	}

	if v.Minor > o.Minor {
		return 1
	} else if v.Minor < o.Minor {
		return -1
	}

	if v.Micro > o.Micro {
		return 1
	} else if v.Micro < o.Micro {
		return -1
	}

	if v.Pre > o.Pre {
		return 1
	} else if v.Pre < o.Pre {
		return -1
	}

	if v.PreNum > o.PreNum {
		return 1
	} else if v.PreNum < o.PreNum {
		return -1
	}

	return 0
}

func (v Version) Equal(o *Version) bool {
	return v.Compare(o) == 0
}

func (v Version) LessThan(o *Version) bool {
	return v.Compare(o) == -1
}

func (v Version) LessThanOrEqual(o *Version) bool {
	return v.Compare(o) <= 0
}

func (v Version) GreaterThan(o *Version) bool {
	return v.Compare(o) == 1
}

func (v Version) GreaterThanOrEqual(o *Version) bool {
	return v.Compare(o) >= 0
}

func (v Version) getPreString() string {
	switch v.Pre {
	case -1:
		return "rc" + strconv.Itoa(v.PreNum)
	case -2:
		return "b" + strconv.Itoa(v.PreNum)
	case -3:
		return "a" + strconv.Itoa(v.PreNum)
	default:
		return ""
	}
}

func (v Version) getStringWithoutPre() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Micro)
}

func (v Version) getFullString() string {
	return v.getStringWithoutPre() + v.getPreString()
}

func (v Version) String() string {
	return v.getFullString()
}

// if ?.0.0 -> 1
// if ?.?.0 -> 2
// if ?.?.? -> 3
// if ?.?.?[a,b,rc]? -> 4
func (v Version) Count() int {
	if v.Pre != 0 {
		return 4
	}
	if v.Micro != 0 {
		return 3
	}
	if v.Minor != 0 {
		return 2
	}
	if v.Major != 0 {
		return 1
	}
	return 0
}

func (v *Version) Set(val string) error {
	var err error
	*v, err = NewVersion(val)
	return err
}

func (v *Version) Type() string {
	return "version"
}

func (v *Version) UnmarshalJSON(data []byte) error {
	var err error
	*v, err = NewVersion(strings.Trim(string(data), `"`))
	return err
}

func (v Version) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, v.getFullString())), nil
}

func PrintVersions(vs []Version) {
	for _, v := range vs {
		fmt.Println(v)
	}
}

func PrintVersionMap(vm map[int]Version) {
	for k, v := range vm {
		fmt.Printf("%d: %+v\n", k, v)
	}
}
