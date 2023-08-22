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
	"golang.org/x/sys/windows/registry"
)

type RegistryInfo struct {
	DisplayName    string
	Version        string
	DirectoryPath  string
	ExecutablePath string
}

type RegistryInfoMap map[string]map[string]RegistryInfo

func readRegistryFrom(from registry.Key) (RegistryInfoMap, error) {
	registryData := make(RegistryInfoMap)

	companyList, err := registry.OpenKey(from, "Software\\Python", registry.READ)
	if err != nil {
		return registryData, err
	}
	defer companyList.Close()

	companyListNames, err := companyList.ReadSubKeyNames(-1)
	if err != nil {
		return registryData, err
	}

	for _, company := range companyListNames {
		// skip PyLauncher. reserved and not company
		if company == "PyLauncher" {
			continue
		}
		// skip not PythonCore. not support other company.
		// TODO: support other company.
		if company != "PythonCore" {
			continue
		}

		tagList, err := registry.OpenKey(companyList, company, registry.READ)
		if err != nil {
			continue
		}
		defer tagList.Close()

		tagListNames, err := tagList.ReadSubKeyNames(-1)
		if err != nil {
			continue
		}

		companyMap := make(map[string]RegistryInfo)
		registryData[company] = companyMap

		for _, tag := range tagListNames {
			tagInfo, err := registry.OpenKey(tagList, tag, registry.READ)
			if err != nil {
				continue
			}
			defer tagInfo.Close()

			var info RegistryInfo
			info.DisplayName, _, err = tagInfo.GetStringValue("DisplayName")
			if err != nil {
				continue
			}
			info.Version, _, err = tagInfo.GetStringValue("Version")
			if err != nil {
				continue
			}

			installPath, err := registry.OpenKey(tagInfo, "InstallPath", registry.READ)
			if err != nil {
				continue
			}
			defer installPath.Close()

			info.DirectoryPath, _, err = installPath.GetStringValue("")
			if err != nil {
				continue
			}

			info.ExecutablePath, _, err = installPath.GetStringValue("ExecutablePath")
			if err != nil {
				continue
			}

			companyMap[tag] = info
		}
	}
	return registryData, nil
}

func readRegistry() (RegistryInfoMap, error) {
	HKLM, _ := readRegistryFrom(registry.LOCAL_MACHINE)
	HKCU, _ := readRegistryFrom(registry.CURRENT_USER)

	registryData := HKLM
	for company, tagInfo := range HKCU {
		for tag, info := range tagInfo {
			if _, ok := registryData[company][tag]; !ok {
				registryData[company][tag] = info
			}
		}
	}
	return registryData, nil
}

func getPythonVersions(config Config) []Version {
	versions := []Version{}

	if registryData, err := readRegistry(); err == nil {
		// TODO: support other company.
		for _, tagInfo := range registryData["PythonCore"] {
			version_string := tagInfo.Version
			v, err := NewVersion(version_string)
			if err != nil {
				continue
			}
			versions = append(versions, v)
		}
	}

	return versions
}
