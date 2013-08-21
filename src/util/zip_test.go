package util_test

import (
	"io/ioutil"
	// "os"
	"path/filepath"
	"testing"
	. "util"
)

func TestUnzip(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	println(tempDir)
	// defer os.RemoveAll(tempDir)
	options := map[string]int{"skip-components": 1}
	Unzip(filepath.Join("testdata", "test.zip"), tempDir, options)
	if content, err := ioutil.ReadFile(filepath.Join(tempDir, "textfile.txt")); err == nil {
		println(string(content))
	} else {
		t.Error(err)
	}
}
