package go2ts

import (
	"fmt"
	"testing"
)

func TestReadTypes(t *testing.T) {
	s, err := ReadTypes("testdata/example")
	if err != nil {
		t.Fatal(err.Error())
	}
	// TODO
	fmt.Println(s)
}
