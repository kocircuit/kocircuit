package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type SourceRepo []*SourceFile

func (sr SourceRepo) Materialize(rootDir string) error {
	for _, sf := range sr {
		if err := sf.Materialize(rootDir); err != nil {
			return fmt.Errorf("materializing %v (%v)", Sprint(sf), err)
		}
	}
	return nil
}

type SourceFile struct {
	Dir  string // "ko/pkg/path/func"
	Base string // "Main.go"
	Body string // source
}

func (sf *SourceFile) Path() string { return path.Join(sf.Dir, sf.Base) }

func (sf *SourceFile) Materialize(rootDir string) error {
	if err := os.MkdirAll(path.Join(rootDir, sf.Dir), 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path.Join(rootDir, sf.Path()), []byte(sf.Body), 0644); err != nil {
		return err
	}
	return nil
}
