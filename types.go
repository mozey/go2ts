package go2ts

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// ReadTypes for all files in packagePath
func ReadTypes(packagePath string) (s string, err error) {
	// List files
	info, err := os.Stat(packagePath)
	if err != nil {
		if os.IsNotExist(err) {
			return s, errors.Errorf("path %s does not exist", packagePath)
		}
	}
	if !info.IsDir() {
		return s, errors.Errorf("path %s is not a directory", packagePath)
	}
	files, err := os.ReadDir(packagePath)
	if err != nil {
		return s, errors.WithStack(err)
	}

	fileSet := token.NewFileSet()
	var convertBuf bytes.Buffer
	addNode := func(filePath, name string, n ast.Node) {
		var nodeBuf bytes.Buffer
		// String representation of ast.Node
		// https://stackoverflow.com/a/52619499/639133
		printer.Fprint(&nodeBuf, fileSet, n)
		convertBuf.WriteString(
			fmt.Sprintf("// %s#%s\n", filePath, name))
		convertBuf.WriteString(fmt.Sprintf("type %s ", name))
		convertBuf.Write(nodeBuf.Bytes())
		convertBuf.WriteString("\n\n")
	}

	for _, file := range files {
		// Read file
		filePath := filepath.Join(packagePath, file.Name())
		info, err := os.Stat(filePath)
		if err != nil {
			return s, errors.WithStack(err)
		}
		if info.IsDir() {
			continue
		}
		b, err := os.ReadFile(filePath)
		if err != nil {
			return s, errors.WithStack(err)
		}

		// Parse
		f, err := parser.ParseFile(
			fileSet, "editor.go", b, parser.SpuriousErrors)
		if err != nil {
			return s, errors.WithStack(err)
		}

		name := ""
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.Ident:
				name = x.Name
			case *ast.ArrayType:
				addNode(filePath, name, n)
				return false
			case *ast.StructType:
				addNode(filePath, name, n)
				return false
			}
			return true
		})
	}

	return convertBuf.String(), nil
}
