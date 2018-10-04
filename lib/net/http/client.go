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
	"net"
	"net/http"
	"sync"
	"time"

	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("NewClient", new(goNewClient))
}

// Client write http.Client
type Client struct {
	http.Client
}

// newClient creates a new Client with default settings.
func newClient() *Client {
	return &Client{
		Client: http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}
}

// getDefaultClient returns the default client, creating one if needed.
func getDefaultClient() *Client {
	defaultClientOnce.Do(func() {
		defaultClient = newClient()
	})
	return defaultClient
}

var (
	defaultClient     *Client
	defaultClientOnce sync.Once
)

type goNewClient struct {
}

func (g *goNewClient) Play(ctx *runtime.Context) *Client {
	return newClient()
}

func (g *goNewClient) Doc() string {
	return "NewClient() creates a new HTTP client"
}

func (g *goNewClient) Help() string {
	return "NewClient()"
}
