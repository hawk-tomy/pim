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
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"time"
)

type VersionCache struct {
	UpdateDate            time.Time       `json:"update_date"`
	AllVersions           []Version       `json:"versions"`
	FailedMinimumVersions map[int]Version `json:"failed_minimum_versions"`
}

var (
	cacheDir         string
	versionCacheFile string
)

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	cacheDir = filepath.Join(home, ".cache", "pim")
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(cacheDir, 0755)
		cobra.CheckErr(err)
	}

	versionCacheFile = filepath.Join(cacheDir, "cache.json")
}

func readCache() (bool, error) {
	var cache VersionCache
	cacheFile, err := os.Open(versionCacheFile)
	if err != nil {
		return true, err
	}
	defer deferErrCheck(cacheFile.Close)

	byteValue, err := io.ReadAll(cacheFile)
	if err != nil {
		return true, err
	}

	err = json.Unmarshal(byteValue, &cache)

	if err != nil {
		return true, err
	}

	// always use cache. maybe add key-value pair, but do not change in cached.
	if v := cache.FailedMinimumVersions; v != nil {
		failedMinimumVersions = v
	}

	if time.Now().Compare(cache.UpdateDate.Add(time.Hour*24)) > 0 {
		return true, nil
	}

	saveVersions(cache.AllVersions)

	return false, nil
}

func saveCache() error {
	cacheFile, err := os.Create(versionCacheFile)
	if err != nil {
		return err
	}
	defer deferErrCheck(cacheFile.Close)

	cache := VersionCache{
		UpdateDate:            time.Now().UTC(),
		AllVersions:           allVersions,
		FailedMinimumVersions: failedMinimumVersions,
	}

	byteValue, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	_, err = cacheFile.Write(byteValue)
	return err
}

func CleanCache() error {
	return os.RemoveAll(cacheDir)
}
