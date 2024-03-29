package templates

import (
	"os"

	"github.com/snemitz/stubborn/pkg/stub/templates/set"
	"github.com/snemitz/stubborn/pkg/stub/templates/vanilla"
)

func NewSet(path string, pkg string) set.Writer {
	return set.Writer{
		Path: path,
		Pkg:  pkg,
	}
}

func NewVanilla(out *os.File, pkg string) vanilla.Writer {
	return vanilla.Writer{
		Out: out,
		Pkg: pkg,
	}
}
