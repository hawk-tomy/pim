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
package cmd

import (
	"os"

	"github.com/hawk-tomy/pim/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	config  lib.Config
)

var rootCmd = &cobra.Command{
	Use:   "pim",
	Short: "Install/Update python.",
	Long: `Install/Update python.

Install/Update python with select/latest version.
You can install, update, show installed version.

Now, only support python/cpython. (PR is welcome!)`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/pim/config.toml)")

	rootCmd.PersistentFlags().BoolP("pre-release", "p", false, "allow pre-release lib")
	viper.BindPFlag("AllowPreRelease", rootCmd.PersistentFlags().Lookup("pre-release"))

	rootCmd.PersistentFlags().BoolP("all-user", "a", false, "install for all user")
	viper.BindPFlag("ForAllUser", rootCmd.PersistentFlags().Lookup("all-user"))

	rootCmd.PersistentFlags().StringP("target-directory", "t", "", "install target directory")
	viper.BindPFlag("TargetDirectory", rootCmd.PersistentFlags().Lookup("target-directory"))

	rootCmd.PersistentFlags().StringToStringVarP(
		&config.AdditionalInstallerOptions,
		"additional-option",
		"o",
		nil,
		`additional install option.

if windows, see https://docs.python.org/3.12/using/windows.html#installing-without-ui
e.g. enable 'CompileAll' =>  "-o CompileAll=1" or "--additional-option=CompileAll=1"
`,
	)
	viper.BindPFlag("AdditionalInstallerOptions", rootCmd.PersistentFlags().Lookup("additional-options"))

	rootCmd.PersistentFlags().BoolVarP(&lib.SkipConfirm, "force", "f", false, "skip confirmation")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile(lib.ConfigPath)
	}

	viper.ReadInConfig()
	if err := viper.Unmarshal(&config); err != nil {
		cobra.CheckErr(err)
	}
}
