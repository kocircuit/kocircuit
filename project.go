package ko

import (
	"strings"
)

func SanitizeKoCompilerSourcePath(file string) string {
	if i := strings.Index(file, "github.com/kocircuit/kocircuit"); i >= 0 {
		return file[i+len("github.com/kocircuit/"):]
	}
	return file
}
