package go2ts

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"

	"github.com/fatih/structtag"
	"github.com/pkg/errors"
)

var Indent = "    "

// TSTypePrefix to use when generating the TypeScript types.
// See discussion re. "declare" vs "export"
// https://stackoverflow.com/q/35019987/639133
var TSTypePrefix = "declare"

func getIdent(s string) string {
	switch s {
	case "bool":
		return "boolean"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64",
		"complex64", "complex128":
		return "number"
	}

	return s
}

func writeType(s *strings.Builder, t ast.Expr, depth int, optionalParens bool) error {
	switch t := t.(type) {
	case *ast.StarExpr:
		if optionalParens {
			s.WriteByte('(')
		}
		err := writeType(s, t.X, depth, false)
		if err != nil {
			return errors.WithStack(err)
		}
		s.WriteString(" | undefined")
		if optionalParens {
			s.WriteByte(')')
		}
	case *ast.ArrayType:
		if v, ok := t.Elt.(*ast.Ident); ok && v.String() == "byte" {
			s.WriteString("string")
			break
		}
		err := writeType(s, t.Elt, depth, true)
		if err != nil {
			return errors.WithStack(err)
		}
		s.WriteString("[]")
	case *ast.StructType:
		s.WriteString("{\n")
		writeFields(s, t.Fields.List, depth+1)

		for i := 0; i < depth+1; i++ {
			s.WriteString(Indent)
		}
		s.WriteByte('}')
	case *ast.Ident:
		s.WriteString(getIdent(t.String()))
	case *ast.SelectorExpr:
		longType := fmt.Sprintf("%s.%s", t.X, t.Sel)
		switch longType {
		case "time.Time":
			s.WriteString("string")
		case "decimal.Decimal":
			s.WriteString("number")
		default:
			s.WriteString(longType)
		}
	case *ast.MapType:
		s.WriteString("{ [key: ")
		err := writeType(s, t.Key, depth, false)
		if err != nil {
			return errors.WithStack(err)
		}
		s.WriteString("]: ")
		err = writeType(s, t.Value, depth, false)
		if err != nil {
			return errors.WithStack(err)
		}
		s.WriteByte('}')
	case *ast.InterfaceType:
		s.WriteString("any")
	default:
		err := fmt.Errorf("unhandled: %s, %T", t, t)
		return errors.WithStack(err)
	}
	return nil
}

var validJSNameRegexp = regexp.MustCompile(`(?m)^[\pL_][\pL\pN_]*$`)

func validJSName(n string) bool {
	return validJSNameRegexp.MatchString(n)
}

func writeFields(s *strings.Builder, fields []*ast.Field, depth int) error {
	for _, f := range fields {
		optional := false

		var fieldName string
		if len(f.Names) != 0 && f.Names[0] != nil && len(f.Names[0].Name) != 0 {
			fieldName = f.Names[0].Name
		}
		if len(fieldName) == 0 || 'A' > fieldName[0] || fieldName[0] > 'Z' {
			continue
		}

		var name string
		if f.Tag != nil {
			tags, err := structtag.Parse(f.Tag.Value[1 : len(f.Tag.Value)-1])
			if err != nil {
				return errors.WithStack(err)
			}

			jsonTag, err := tags.Get("json")
			if err == nil {
				name = jsonTag.Name
				if name == "-" {
					continue
				}

				optional = jsonTag.HasOption("omitempty")
			}
		}

		if len(name) == 0 {
			name = fieldName
		}

		for i := 0; i < depth+1; i++ {
			s.WriteString(Indent)
		}

		quoted := !validJSName(name)

		if quoted {
			s.WriteByte('\'')
		}
		s.WriteString(name)
		if quoted {
			s.WriteByte('\'')
		}

		switch t := f.Type.(type) {
		case *ast.StarExpr:
			optional = true
			f.Type = t.X
		}

		if optional {
			s.WriteByte('?')
		}

		s.WriteString(": ")

		writeType(s, f.Type, depth, false)

		s.WriteString(";\n")
	}
	return nil
}

const wrapper = `package main

func main() {
	%s
}`

func Convert(s string) (ts string, err error) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s, nil
	}

	fileSet := token.NewFileSet()
	var f ast.Node
	f, err = parser.ParseExprFrom(fileSet, "editor.go", s, parser.SpuriousErrors)
	if err != nil {
		s = fmt.Sprintf(wrapper, s)

		f, err = parser.ParseFile(fileSet, "editor.go", s, parser.SpuriousErrors)
		if err != nil {
			return ts, errors.WithStack(err)
		}
	}

	w := new(strings.Builder)
	name := "MyInterface"

	first := true

	var builderErr error
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			name = x.Name

		// TODO If Go type declaration is preceded by comment lines,
		// then preserve the comment in the TypeScript declaration.
		// See examples in testdata/example/compare/ReadTypes.txt

		case *ast.ArrayType:
			if !first {
				w.WriteString("\n\n")
			}

			w.WriteString(fmt.Sprintf("%s interface ", TSTypePrefix))
			w.WriteString(name)
			// How can I define an interface for an array of objects?
			// https://stackoverflow.com/a/25470775/639133
			w.WriteString(" extends Array<")
			w.WriteString(fmt.Sprintf("%s", x.Elt))
			w.WriteString(">{}")
			return false

		case *ast.StructType:
			if !first {
				w.WriteString("\n\n")
			}

			w.WriteString(fmt.Sprintf("%s interface ", TSTypePrefix))
			w.WriteString(name)
			w.WriteString(" {\n")

			err = writeFields(w, x.Fields.List, 0)
			if err != nil {
				builderErr = err
				return false
			}

			w.WriteByte('}')

			first = false

			return false
		}

		return true
	})
	if builderErr != nil {
		return ts, builderErr
	}

	return w.String(), nil
}
