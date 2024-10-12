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
	"errors"
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

			needLatest, err := cmd.Flags().GetBool("latest")
			if err != nil {
				return err
			}

			versionInfo := version.String()
			if needLatest {
				versionInfo = fmt.Sprintf("%d.%d.x (detect latest)", version.Major, version.Minor)
			}

			fmt.Printf("install options\n")
			fmt.Printf("  version: %s\n", versionInfo)
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
				err = lib.InstallPython(config, version, needLatest)
				if err != nil {
					var sErr *lib.StatusError
					if errors.As(err, &sErr) {
						fmt.Printf("failed to download installer. version: %s, status: %d\n", version.String(), sErr.Status)
						return nil
					}
					return err
				}

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
