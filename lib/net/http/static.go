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

package http

import (
	"net/http"

	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("ServeLocalDir", new(goServeLocalDir))
}

type goServeLocalDir struct {
	// The desired http server address
	Address string `ko:"name=address"`
	// Local directory
	Dir string `ko:"name=dir"`
	// Root URL for serving dir
	URL string `ko:"name=url"` // root URL
}

func (g *goServeLocalDir) Play(ctx *runtime.Context) bool {
	http.Handle(g.URL, http.FileServer(http.Dir(g.Dir)))
	ctx.Printf("serving %q on %s%s", g.Dir, g.Address, g.URL)
	if err := http.ListenAndServe(g.Address, nil); err != nil {
		ctx.Fatalf("serving %q on %s%s (%v)", g.Dir, g.Address, g.URL, err)
		return false // never reached
	}
	return true
}

func (g *goServeLocalDir) Doc() string {
	return "ServeLocalDir(address, dir, url) runs an HTTP server that serves files from a local directory"
}

func (g *goServeLocalDir) Help() string {
	return "ServeLocalDir(address, dir, url)"
}
