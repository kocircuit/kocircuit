package syntax

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Repository interface {
	Read(filePath string) (content string, err error)
	List(dirPath string) (file, subdir []string, err error)
}

func NewLocalRepositories(rootDirs []string) Repository {
	r := make(repositories, len(rootDirs))
	for i, dir := range rootDirs {
		r[i] = NewLocalRepository([]string{dir})
	}
	return r
}

type repositories []Repository

func (r repositories) Read(filePath string) (content string, err error) {
	for _, r := range r {
		if content, err = r.Read(filePath); err == nil {
			return content, nil
		}
	}
	return "", fmt.Errorf("file %q not found in any repository", filePath)
}

func (r repositories) List(dirPath string) (file, subdir []string, err error) {
	for _, r := range r {
		if file, subdir, err = r.List(dirPath); err == nil {
			return file, subdir, nil
		}
	}
	return nil, nil, fmt.Errorf("directory %q not found in any repository", dirPath)
}

func NewLocalRepository(rootDirs []string) Repository {
	return &localRepository{roots: rootDirs}
}

type localRepository struct {
	roots []string
}

func (repo *localRepository) Read(filePath string) (string, error) {
	var firstErr error
	for _, root := range repo.roots {
		buf, err := ioutil.ReadFile(path.Join(root, filePath))
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		return string(buf), nil
	}
	return "", firstErr
}

func (repo *localRepository) List(dirPath string) (file, subdir []string, err error) {
	for _, root := range repo.roots {
		d, err := os.Open(path.Join(root, dirPath))
		if err != nil {
			continue
		}
		defer d.Close()
		ff, err := d.Readdir(0)
		if err != nil {
			return nil, nil, err
		}
		for _, f := range ff {
			if f.IsDir() {
				subdir = append(subdir, path.Join(dirPath, f.Name()))
			} else {
				if path.Ext(f.Name()) == ".ko" { // so that GOPATH = KOPATH is ok
					file = append(file, path.Join(dirPath, f.Name()))
				}
			}
		}
		return file, subdir, nil
	}
	return nil, nil, nil
}

// SrcDir is repository.
type SrcDir map[string]interface{} // name -> SrcDir or SrcFile

// SrcFile is a source file in a SrcDir repository.
type SrcFile string

func (srcDir SrcDir) Read(filePath string) (string, error) {
	dir, base := path.Split(filePath)
	for _, k := range splitPath(dir) {
		if subdir, ok := srcDir[k]; !ok {
			return "", fmt.Errorf("path %s not found", filePath)
		} else if sd, ok := subdir.(SrcDir); !ok {
			return "", fmt.Errorf("path %s not found", filePath)
		} else {
			srcDir = sd
		}
	}
	switch u := srcDir[base].(type) {
	case SrcDir:
		return "", fmt.Errorf("path %s is a directory", filePath)
	case SrcFile:
		return string(u), nil
	}
	return "", fmt.Errorf("file %s not found", filePath)
}

func (srcDir SrcDir) List(dirPath string) (file, subdir []string, err error) {
	for k, v := range srcDir {
		switch v.(type) {
		case SrcFile:
			file = append(file, k)
		case SrcDir:
			subdir = append(subdir, k)
		default:
			panic("not a file or directory")
		}
	}
	return file, subdir, nil
}

func splitPath(repoPath string) []string {
	repoPath = path.Clean(repoPath)
	repoPath = strings.TrimLeft(repoPath, "/")
	repoPath = strings.TrimRight(repoPath, "/")
	return strings.Split(repoPath, "/")
}
