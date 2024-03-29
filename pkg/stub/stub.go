package stub

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

type visitor struct {
	level  int
	result []InterfaceDecl
}

type InterfaceDecl struct {
	Name    string
	Methods []MethodDecl
}

type MethodDecl struct {
	Name   string
	Params []ParamDecl
	Return []ParamDecl
}

type ParamDecl struct {
	Name string
	Type string
}

func parseParam(field *ast.Field) (ParamDecl, error) {
	var pType []string
	var name string
	switch paramType := field.Type.(type) {
	case *ast.ArrayType:
		ident, isIdent := paramType.Elt.(*ast.Ident)
		if !isIdent {
			return ParamDecl{}, ErrUnknownProperty
		}
		pType = []string{"[]" + ident.Name}
	case *ast.Ident:
		pType = []string{paramType.Name}
	case *ast.SelectorExpr:
		ident, isIdent := paramType.X.(*ast.Ident)
		if !isIdent {
			return ParamDecl{}, ErrUnknownProperty
		}
		pType = []string{
			ident.Name,
			paramType.Sel.Name,
		}
	}
	if len(field.Names) > 0 {
		name = field.Names[0].Name
	}

	return ParamDecl{
		Name: name,
		Type: strings.Join(
			pType,
			".",
		),
	}, nil
}

func parseInterface(node *ast.InterfaceType) ([]MethodDecl, error) {
	funcs := []MethodDecl{}
	if node.Incomplete {
		return nil, ErrIncompleteCode
	}
	for _, method := range node.Methods.List {
		funcType, ok := method.Type.(*ast.FuncType)
		if !ok {
			continue
		}
		m := MethodDecl{
			Name:   method.Names[0].Name,
			Params: []ParamDecl{},
			Return: []ParamDecl{},
		}
		for _, parameter := range funcType.Params.List {
			param, err := parseParam(parameter)
			if err != nil {
				continue
			}
			m.Params = append(m.Params, param)
		}
		for _, result := range funcType.Results.List {
			param, err := parseParam(result)
			if err != nil {
				continue
			}
			m.Return = append(m.Return, param)
		}
		funcs = append(funcs, m)
	}
	return funcs, nil
}

func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch d := n.(type) {
	case *ast.TypeSpec:
		switch typeSpec := d.Type.(type) {
		case *ast.InterfaceType:
			fns, err := parseInterface(typeSpec)
			if err != nil {
				return nil
			}
			r := InterfaceDecl{
				Name:    d.Name.Name,
				Methods: fns,
			}
			v.result = append(v.result, r)
		}
		return nil
	case *ast.InterfaceType:
		return nil
	case *ast.File:
		// name of package
	default:
		v.level++
	}
	return v
}

func Get(w Writer) error {
	f, err := parser.ParseFile(token.NewFileSet(), "../example/basic/interface.go", nil, 0)
	if err != nil {
		return errors.Wrap(err, "parsing input file")
	}
	visitor := &visitor{
		level:  0,
		result: []InterfaceDecl{},
	}
	ast.Walk(visitor, f)

	err = w.Write(visitor.result[0])

	return errors.Wrap(err, "writing interface stub")
}
