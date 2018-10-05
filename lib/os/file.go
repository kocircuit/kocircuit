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

package os

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

type FileSystem struct {
	Root *Dir `ko:"name=root"` // local root
}

type Dir struct {
	Name string  `ko:"name=name"`
	Dir  []*Dir  `ko:"name=dir"`
	File []*File `ko:"name=file"`
}

type File struct {
	Name string `ko:"name=name"`
	Body string `ko:"name=body"`
}

func NewDir(name string) *Dir {
	return &Dir{Name: name}
}

func (dir *Dir) AddFile(name, body string) {
	dir.File = append(dir.File, &File{Name: name, Body: body})
}

func (dir *Dir) Walk(path []string) *Dir {
	for _, p := range path {
		dir = dir.walk(p)
	}
	return dir
}

func (dir *Dir) walk(name string) *Dir {
	for _, d := range dir.Dir {
		if d.Name == name {
			return d
		}
	}
	d := &Dir{Name: name}
	dir.Dir = append(dir.Dir, d)
	return d
}

func (dir *Dir) Materialize(atPath string) error {
	q := []*pendingDir{{PathTo: atPath, Dir: dir}}
	for len(q) > 0 {
		p := path.Join(q[0].PathTo, q[0].Dir.Name)
		if err := os.MkdirAll(p, 0755); err != nil {
			return err
		}
		for _, f := range q[0].Dir.File {
			if err := ioutil.WriteFile(path.Join(p, f.Name), []byte(f.Body), 0644); err != nil {
				return err
			}
		}
		for _, d := range q[0].Dir.Dir {
			q = append(q, &pendingDir{PathTo: p, Dir: d})
		}
		q = q[1:]
	}
	return nil
}

type pendingDir struct {
	PathTo string
	Dir    *Dir
}

func GetNewTempDir() string {
	cd := path.Join(os.TempDir(), fmt.Sprintf("ko/os-%s", runtime.ExecutionID()))
	if err := os.MkdirAll(cd, 0755); err != nil {
		panic(err)
	}
	return cd
}
