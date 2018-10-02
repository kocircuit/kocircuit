//
// Copyright © 2018 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package syntax

import (
	"fmt"
	"sort"
)

// ParseRepo returns a map from package path to a set of parsed files from the package directory.
func ParseRepo(repo Repository, pkgPaths []string) (map[string][]File, error) {
	pending := make(map[string]bool) // pending is a set of paths to unparsed package directories
	for _, pkgPath := range pkgPaths {
		pending[pkgPath] = true
	}
	parsed := map[string]bool{} // parsed is a set of paths to parsed packages
	prog := map[string][]File{} // package path -> parsed package files
	for len(pending) > 0 {
		order := make([]string, 0, len(pending))
		for pkg := range pending {
			if !parsed[pkg] {
				order = append(order, pkg)
			}
		}
		pending = map[string]bool{}
		sort.Strings(order)
		for _, pkg := range order {
			pkgFiles, err := parsePackage(repo, pkg)
			if err != nil {
				return nil, err
			}
			prog[pkg], parsed[pkg] = pkgFiles, true
			for _, file := range pkgFiles {
				for _, imp := range file.Import {
					if _, parsed := prog[imp.Path]; !parsed {
						pending[imp.Path] = true
					}
				}
			}
		}
	}
	return prog, nil
}

func parsePackage(repo Repository, pkgPath string) (file []File, err error) {
	sourceFile, _, err := repo.List(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("listing package directory %q (%v)", pkgPath, err)
	}
	for _, sf := range sourceFile {
		text, err := repo.Read(sf)
		if err != nil {
			return nil, fmt.Errorf("reading source file %q (%v)", sf, err)
		}
		f, err := ParseFileString(sf, text)
		if err != nil {
			return nil, err
		}
		file = append(file, f)
	}
	return file, nil
}
