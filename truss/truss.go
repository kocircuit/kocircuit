// Package truss provides a framework for source manipulations within a repo.
package truss

type Truss struct {
	Repo string       `ko:"name=repo"`
	Go   *GoToolchain `ko:"name=go"`
}

type GoToolchain struct {
	GoBinaryPath `ko:"name=goBinaryPath"`
}
