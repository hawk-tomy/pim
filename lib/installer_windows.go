//go:build windows

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
	"github.com/hawk-tomy/pim/lib/list"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// example (python 3.11.0)
// arch: amd64/arm64
// - amd64: https://www.python.org/ftp/python/3.11.0/python-3.11.0-amd64.exe
// - arm64: https://www.python.org/ftp/python/3.11.0/python-3.11.0-arm64.exe
// pre-releases: a/b/rc/(final)
// - alpha: https://www.python.org/ftp/python/3.11.0/python-3.11.0a1-amd64.exe
// - beta: https://www.python.org/ftp/python/3.11.0/python-3.11.0b1-amd64.exe
// - rc: https://www.python.org/ftp/python/3.11.0/python-3.11.0rc1-amd64.exe
// - final: https://www.python.org/ftp/python/3.11.0/python-3.11.0-amd64.exe
const (
	downloadUrlBase = `https://www.python.org/ftp/python/%s/python-%s-%s.exe`
	fileNameBase    = `python-%s-%s.exe`
)

var (
	installerCacheDir string
)

type StatusError struct {
	Status int
}

func (e *StatusError) Error() string { return fmt.Sprintf("Bad Status: %d", e.Status) }

func init() {
	installerCacheDir = filepath.Join(cacheDir, "installer")
	if _, err := os.Stat(installerCacheDir); os.IsNotExist(err) {
		err := os.MkdirAll(installerCacheDir, 0755)
		cobra.CheckErr(err)
	}
}

func downloadInstaller(version Version) (string, error) {
	dirVersionString := version.getStringWithoutPre()
	fileVersionString := version.getFullString()
	arch := "amd64"
	if runtime.GOARCH == "arm64" {
		arch = "arm64"
	}

	url := fmt.Sprintf(downloadUrlBase, dirVersionString, fileVersionString, arch)
	filePath := filepath.Join(installerCacheDir, fmt.Sprintf(fileNameBase, fileVersionString, arch))

	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		if v, ok := failedMinimumVersions[version.Minor]; !ok || v.GreaterThan(version) {
			failedMinimumVersions[version.Minor] = version
			cobra.CheckErr(saveCache())
		}
		return "", &StatusError{resp.StatusCode}
	}
	defer deferErrCheck(resp.Body.Close)

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer deferErrCheck(out.Close)

	_, err = io.Copy(out, resp.Body)
	return filePath, err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func buildInstallerArgument(config Config, additionalCmd ...string) []string {
	var args = []string{
		"/quiet",
	}

	for _, cmd := range additionalCmd {
		args = append(args, cmd)
	}

	args = append(args, fmt.Sprintf("InstallAllUsers=%d", boolToInt(config.ForAllUser)))

	if config.TargetDirectory != "" {
		args = append(args, fmt.Sprintf("TargetDir=%s", config.TargetDirectory))
	}

	for k, v := range config.AdditionalInstallerOptions {
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	return args
}

func callInstaller(path string, args ...string) error {
	if WithVerbose > 0 {
		fmt.Printf("call: %s %s\n", path, strings.Join(args, " "))
	}
	cmd := exec.Command(path, args...)

	var stdout strings.Builder
	cmd.Stdout = &stdout
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to call installer: %w\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}
	return nil
}

func doInstall(config Config, version *list.Element[Version], needLatest bool) error {
	var path string
	var err error
	for {
		path, err = downloadInstaller(version.Value)
		if err == nil {
			break
		}
		var sErr *StatusError
		if needLatest && errors.As(err, &sErr) {
			version = version.Prev()
			if version == nil {
				return errors.New("can not found installable version")
			}
			continue
		}
		return err
	}
	return callInstaller(path, buildInstallerArgument(config)...)
}

func doUpdate(config Config, version *list.Element[Version]) error {
	var path string
	var err error
	for {
		path, err = downloadInstaller(version.Value)
		if err == nil {
			break
		}
		var sErr *StatusError
		if errors.As(err, &sErr) {
			version = version.Prev()
			if version == nil {
				return errors.New("can not found installable version. (not found installable version in checked version)")
			}
			if v, ok := installedPythonVersions[version.Value.Minor]; !ok || v.GreaterThanOrEqual(version.Value) {
				// not ok -> UNREACHABLE (check for assert) -> return error
				// v >= version -> prev version is same as or older than already installed version.
				return errors.New("can not found installable version. (not found installable version for newer then installed one.)")
			}
			continue
		}
		return err
	}
	return callInstaller(path, buildInstallerArgument(config)...)
}

func doUninstall(config Config, version Version) error {
	path, err := downloadInstaller(version)
	if err != nil {
		return err
	}
	return callInstaller(path, buildInstallerArgument(config, "/uninstall")...)
}
