package tools

import (
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func Evaluate(src string) string {
	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)

	v, err := i.Eval(src)
	if err != nil {
		panic(err)
	}

	v, err = i.Eval("tool.Do")
	if err != nil {
		panic(err)
	}

	f := v.Interface().(func() string)

	r := f()
	println("r is :", r)
	return r
}
