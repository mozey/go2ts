package go2ts

import (
	"os"
	"testing"
)

func TestConvert(t *testing.T) {
	in, err := os.ReadFile("testdata/example.go.txt")
	if err != nil {
		t.Fatal(err.Error())
	}
	compareOut, err := os.ReadFile("testdata/example.ts.txt")
	if err != nil {
		t.Fatal(err.Error())
	}
	out := Convert(string(in))
	if out != string(compareOut) {
		t.Fatal("Output does not match test data")
	}
}
