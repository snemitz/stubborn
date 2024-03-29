package main

import (
	"github.com/snemitz/stubborn/pkg/stub"
	"github.com/snemitz/stubborn/pkg/stub/templates"
)

func main() {
	var err error
	w := templates.NewSet("../example/basic/impl", "impl")

	err = stub.Get(w)
	if err != nil {
		panic(err)
	}
}
