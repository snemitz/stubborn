package vanilla

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/snemitz/stubborn/pkg/stub"
)

type MethodTemplate struct {
	MethodName string
	Name       string
	Short      string
	Params     string
	Returns    string
}

type StructTemplate struct {
	Pkg   string
	Name  string
	Funcs []MethodTemplate
}

//go:embed *
var f embed.FS

type Writer struct {
	Out *os.File
	Pkg string
}

func (v Writer) Write(i stub.InterfaceDecl) error {
	tmpl, err := template.ParseFS(f, "struct.tmpl")
	if err != nil {
		return errors.Wrap(err, "parsing tmpl")
	}

	lc := strcase.ToDelimited(i.Name, '.')

	splits := strings.Split(lc, ".")
	var short string
	for _, s := range splits {
		short = short + fmt.Sprintf("%c", s[0])
	}
	funcs := []MethodTemplate{}
	for _, fn := range i.Methods {
		params := []string{}
		for _, x := range fn.Params {
			if x.Name != "" {
				params = append(params, x.Name+" "+x.Type)
			} else {
				params = append(params, x.Type)
			}
		}
		returns := []string{}
		for _, x := range fn.Return {
			if x.Name != "" {
				returns = append(returns, x.Name+" "+x.Type)
			} else {
				returns = append(returns, x.Type)
			}
		}
		returnString := strings.Join(returns, ", ")
		if len(returns) > 1 {
			returnString = "(" + returnString + ")"
		}
		funcs = append(funcs, MethodTemplate{
			MethodName: fn.Name,
			Name:       i.Name,
			Short:      short,
			Params:     strings.Join(params, ", "),
			Returns:    returnString,
		})
	}
	err = tmpl.Execute(v.Out, StructTemplate{
		Pkg:   v.Pkg,
		Name:  i.Name,
		Funcs: funcs,
	})
	if err != nil {
		return errors.Wrap(err, "executing tmpl")
	}
	return nil
}
