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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	BASE_URL                     = "https://api.github.com/repos/python/cpython/tags?per_page=100&page=%d"
	supportedMinimumMinorVersion = 9
)

var (
	finalBinaryRelease    map[int]Version
	latestVersions        map[int]Version
	latestVersionsWithPre map[int]Version

	fetchedLatestVersions bool
	maxMinor              int
)

func init() {
	finalBinaryRelease = make(map[int]Version)
	finalBinaryRelease[9] = Version{Major: 3, Minor: 9, Micro: 13, Pre: 0, PreNum: 0}
	finalBinaryRelease[10] = Version{Major: 3, Minor: 10, Micro: 11, Pre: 0, PreNum: 0}

	latestVersions = make(map[int]Version)
	latestVersionsWithPre = make(map[int]Version)
	fetchedLatestVersions = false
	maxMinor = -1
}

type Tag struct {
	Name   string `json:"name"`
	Commit struct {
		Sha string `json:"sha"`
		Url string `json:"url"`
	} `json:"commit"`
	ZipballUrl string `json:"zipball_url"`
	TarballUrl string `json:"tarball_url"`
	NodeId     string `json:"node_id"`
}

func updateMaxMinor(vs map[int]Version) {
	for _, v := range vs {
		maxMinor = max(maxMinor, v.Minor)
	}
}

func getVersionsByPage(config Config, page int) ([]Version, error) {
	url := fmt.Sprintf(BASE_URL, page)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("accept", "application/vnd.github+json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed call API")
	}
	byteArray, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tags []Tag
	err = json.Unmarshal(byteArray, &tags)
	if err != nil {
		return nil, err
	}

	var versions []Version
	for _, tag := range tags {
		// tag Name prefix is not v3, skip
		// has prefix
		if !strings.HasPrefix(tag.Name, "v3") {
			continue
		}
		v, err := NewVersion(tag.Name)
		if err != nil {
			continue
		}
		versions = append(versions, v)
	}

	return versions, nil
}

func filterLatestVersion(config Config) {
	for k, v := range latestVersionsWithPre {
		if k < supportedMinimumMinorVersion {
			continue
		}

		if config.AllowPreRelease || v.Pre == 0 {
			latestVersions[k] = v
		}
	}
}

func innerFetchLatestVersions(config Config) error {
	i := 1
	for {
		versions, err := getVersionsByPage(config, i)
		if err != nil {
			return err
		}

		for _, ver := range versions {
			if v, ok := finalBinaryRelease[ver.Minor]; ok && v.LessThan(&ver) {
				continue
			}

			if v, ok := latestVersionsWithPre[ver.Minor]; !ok || v.LessThan(&ver) {
				latestVersionsWithPre[ver.Minor] = ver
			}
		}

		// GitHub API return tags sorted by name.
		updateMaxMinor(latestVersionsWithPre)
		flag := false
		for i := supportedMinimumMinorVersion; i <= maxMinor; i++ {
			if _, ok := latestVersionsWithPre[i]; !ok {
				flag = true
			}
		}
		if !flag {
			break
		}
		i++
	}

	return nil
}

func fetchLatestVersions(config Config) error {
	if fetchedLatestVersions {
		return nil
	}

	need, err := readCache()
	if need {
		if err := innerFetchLatestVersions(config); err != nil {
			return err
		}

		if err := saveCache(); err != nil {
			return err
		}
	}

	filterLatestVersion(config)
	fetchedLatestVersions = true
	return err
}
