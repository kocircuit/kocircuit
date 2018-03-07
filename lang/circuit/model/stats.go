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
