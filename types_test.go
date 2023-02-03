package go2ts

import (
	"fmt"
	"os"
	"testing"
)

func TestReadTypes(t *testing.T) {
	s, err := ReadTypes("testdata/example")
	if err != nil {
		t.Fatal(err.Error())
	}
	filePath := "testdata/example/compare/ReadTypes.txt"
	compareOut, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err.Error())
	}
	if s != string(compareOut) {
		fmt.Println(s)
		t.Fatalf("Output does not match test data in %s", filePath)
	}

	out := Convert(s)
	filePath = "testdata/example/compare/Convert.txt"
	compareOut, err = os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err.Error())
	}
	if out != string(compareOut) {
		fmt.Println(out)
		t.Fatalf("Output does not match test data in %s", filePath)
	}
}
