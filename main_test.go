package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	fmt.Printf("%#v\n", os.Getenv("TRAVIS_GO_VERSION"))
	t.Fail()
}
