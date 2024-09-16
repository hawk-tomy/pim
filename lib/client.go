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
	"github.com/hawk-tomy/pim/lib/list"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

const (
	BaseUrl                      = "https://api.github.com/repos/python/cpython/tags?per_page=100&page=%d"
	supportedMinimumMinorVersion = 9
)

var (
	allVersions           []Version
	fetchedVersions       map[int]*list.List[Version]
	failedMinimumVersions map[int]Version
	fetchedLatestVersions bool // guard against fetching more than once.
)

func init() {
	fetchedVersions = make(map[int]*list.List[Version])
	failedMinimumVersions = make(map[int]Version)
	fetchedLatestVersions = false
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

func getVersionsByPage(config Config, page int) ([]Version, error) {
	url := fmt.Sprintf(BaseUrl, page)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("accept", "application/vnd.github+json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer deferErrCheck(resp.Body.Close)
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

func innerFetchLatestVersions(config Config) error {
	i := 1
	fetchedVersions_ := make([]Version, 0)
	isFetched := make(map[int]bool)
	maxMinor := -1

	for {
		versions, err := getVersionsByPage(config, i)
		if err != nil {
			return err
		}

		fetchedVersions_ = append(fetchedVersions_, versions...)

		// check if all minor versions from supportedMinimumMinorVersion to latestMinorVersion are fetched.
		// GitHub API return tags sorted by name.
		for _, v := range versions {
			isFetched[v.Minor] = true
			maxMinor = max(maxMinor, v.Minor)
		}
		flag := false
		for i := supportedMinimumMinorVersion - 1; i <= maxMinor; i++ {
			if _, ok := isFetched[i]; !ok {
				flag = true
			}
		}
		if !flag { // ok
			break
		}

		// read next page
		i++
	}

	saveVersions(fetchedVersions_)
	return nil
}

func saveVersions(versions []Version) {
	sort.Slice(versions, func(i, j int) bool { return versions[i].LessThan(versions[j]) })
	allVersions = versions
	fetchVersionsEachMinor := make(map[int][]Version)
	for _, v := range versions {
		if _, ok := fetchVersionsEachMinor[v.Minor]; !ok {
			fetchVersionsEachMinor[v.Minor] = make([]Version, 0)
		}
		fetchVersionsEachMinor[v.Minor] = append(fetchVersionsEachMinor[v.Minor], v)
	}
	for k, v := range fetchVersionsEachMinor {
		if k < supportedMinimumMinorVersion {
			continue
		}
		versionList := list.New[Version]()
		for _, ver := range v {
			versionList.PushBack(ver)
		}
		fetchedVersions[k] = versionList
	}
}

func fetchLatestVersions(config Config) error {
	if fetchedLatestVersions {
		return nil
	}

	if need, err := readCache(); need { // if "need" is true, do not use cache. otherwise, not error.
		if err := innerFetchLatestVersions(config); err != nil {
			return err
		}

		if err := saveCache(); err != nil {
			return err
		}
	} else if err != nil {
		log.Fatal(err) // UNREACHABLE
	}

	fetchedLatestVersions = true
	return nil
}
