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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Config struct {
	AllowPreRelease            bool
	ForAllUser                 bool
	TargetDirectory            string
	AdditionalInstallerOptions map[string]string
}

var (
	configDir  string
	ConfigPath string
)

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	configDir = filepath.Join(home, ".config", "pim")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0755)
		cobra.CheckErr(err)
	}

	ConfigPath = filepath.Join(configDir, "config.toml")
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		f, err := os.Create(ConfigPath)
		cobra.CheckErr(err)
		defer f.Close()

		_, err = f.WriteString("")
		cobra.CheckErr(err)
	}
}
