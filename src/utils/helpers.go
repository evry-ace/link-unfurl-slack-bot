package utils

import (
	"fmt"
	"io/ioutil"
)

const (
	// testdataDir is the path to the testdata directory
	testdataDir = "../../testdata"
)

// ReadTestdataFile reads a file from the testdata directory
func ReadTestdataFile(file string) string {
	f := fmt.Sprintf("%s/%s", testdataDir, file)
	return readFileFromPath(f)
}

// readFileFromPath reads a file from the given path
func readFileFromPath(file string) string {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(f)
}
