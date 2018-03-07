package syntax

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

type Repository interface {
	Read(filePath string) (content string, err error)
	List(dirPath string) (file, subdir []string, err error)
	NotFound() []string
}

func NewLocalRepository(rootDir string) Repository {
	return &localRepository{root: rootDir}
}

type localRepository struct {
	root string
	sync.Mutex
	notFound []string
}

func (repo *localRepository) Read(filePath string) (string, error) {
	buf, err := ioutil.ReadFile(path.Join(repo.root, filePath))
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (repo *localRepository) NotFound() []string {
	return repo.notFound
}

func (repo *localRepository) addNotFound(dirPath string) {
	repo.Lock()
	defer repo.Unlock()
	repo.notFound = append(repo.notFound, dirPath)
}

func (repo *localRepository) List(dirPath string) (file, subdir []string, err error) {
	d, err := os.Open(path.Join(repo.root, dirPath))
	if err != nil {
		repo.addNotFound(dirPath)
		return nil, nil, nil
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
