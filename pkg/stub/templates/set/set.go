package set

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/snemitz/stubborn/pkg/stub"
)

type MethodTemplate struct {
	Pkg        string
	MethodName string
	Name       string
	Short      string
	Params     string
	Returns    string
}

type StructTemplate struct {
	Pkg  string
	Name string
}

//go:embed *
var f embed.FS

type Writer struct {
	Path string
	Pkg  string
}

func (s Writer) Write(i stub.InterfaceDecl) error {
	if _, err := os.Stat(s.Path); os.IsNotExist(err) {
		err := os.Mkdir(s.Path, 0755)
		if err != nil {
			return errors.Wrap(err, "creating dir")
		}
	}

	tmplStruct, err := template.ParseFS(f, "struct.tmpl")
	if err != nil {
		return errors.Wrap(err, "parsing tmpl")
	}

	tmplFunc, err := template.ParseFS(f, "func.tmpl")
	if err != nil {
		return errors.Wrap(err, "parsing tmpl")
	}

	lc := strcase.ToDelimited(i.Name, '.')

	splits := strings.Split(lc, ".")
	var short string
	for _, s := range splits {
		short = short + fmt.Sprintf("%c", s[0])
	}
	out, err := os.Create(filepath.Join(s.Path, "struct.go"))
	if err != nil {
		return errors.Wrap(err, "opening struct file")
	}

	err = tmplStruct.Execute(out, StructTemplate{
		Pkg:  s.Pkg,
		Name: i.Name,
	})
	if err != nil {
		return errors.Wrap(err, "executing tmpl")
	}

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

		filename := strcase.ToSnake(fn.Name)
		out, err := os.Create(filepath.Join(s.Path, filename+".go"))
		if err != nil {
			return errors.Wrap(err, "opening struct file")
		}

		err = tmplFunc.Execute(out, MethodTemplate{
			Pkg:        s.Pkg,
			MethodName: fn.Name,
			Name:       i.Name,
			Short:      short,
			Params:     strings.Join(params, ", "),
			Returns:    returnString,
		})
		if err != nil {
			return errors.Wrap(err, "writing func file")
		}
	}
	return nil
}
