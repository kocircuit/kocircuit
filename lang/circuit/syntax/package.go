package syntax

import (
	"fmt"
	"sort"
)

// ParseRepo returns a map from package path to a set of parsed files from the package direcrtory.
func ParseRepo(repo Repository, pkgPath string) (map[string][]File, error) {
	pending := map[string]bool{pkgPath: true} // pending is a set of paths to unparsed package directories
	parsed := map[string]bool{}               // parsed is a set of paths to parsed packages
	prog := map[string][]File{}               // package path -> parsed package files
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
