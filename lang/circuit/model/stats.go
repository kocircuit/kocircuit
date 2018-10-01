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

package model

type RepoStats struct {
	TotalFunc   int64
	TotalStep   int64
	StepPerFunc float64
}

func (repo Repo) Stats() *RepoStats {
	stats := &RepoStats{}
	for _, pkg := range repo {
		for _, f := range pkg {
			stats.TotalFunc += 1
			stats.TotalStep += int64(len(f.Step) - len(f.Field) - 1) // enter/select steps for args accounted for by respective leave steps
		}
	}
	if stats.TotalFunc > 0 {
		stats.StepPerFunc = float64(stats.TotalStep) / float64(stats.TotalFunc)
	}
	return stats
}
