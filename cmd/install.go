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
	"fmt"

	"github.com/hawk-tomy/pim/lib"
	"github.com/spf13/cobra"
)

var (
	installCmd = &cobra.Command{
		Use:   "install",
		Short: "install python",
		Long:  "install python",

		Args: cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				fmt.Println("1st argument must be version.")
			}

			version, err := lib.NewVersion(args[0])
			if err != nil {
				return err
			}

			if need, err := cmd.Flags().GetBool("latest"); err != nil {
				return err
			} else if need {
				version, err = lib.GetLatestVersion(config, version)
				if err != nil {
					return err
				}
			}

			fmt.Printf("install options\n")
			fmt.Printf("  version: %s\n", version.String())
			fmt.Printf("  for all user: %s\n", lib.YesOrNo(config.ForAllUser))
			if config.TargetDirectory != "" {
				fmt.Printf("  install path: %s\n", config.TargetDirectory)
			} else if config.ForAllUser {
				fmt.Printf("  install path: default(all user)\n")
			} else {
				fmt.Printf("  install path: default(only you)\n")
			}

			if lib.Confirm("continue? [Y/n]: ") {
				fmt.Printf("installing...\n")
				return lib.InstallPython(config, version)
			} else {
				fmt.Println("canceled.")
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().BoolP("latest", "l", false, "install latest version (each minor version)")
}
