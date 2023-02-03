package go2ts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// ReadTypes for all files in packagePath
func ReadTypes(packagePath string) (s string, err error) {
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
	for _, file := range files {
		filePath := filepath.Join(packagePath, file.Name())
		b, err := os.ReadFile(filePath)
		if err != nil {
			return s, errors.WithStack(err)
		}
		// TODO
		fmt.Println(b)
	}
	return s, nil
}
