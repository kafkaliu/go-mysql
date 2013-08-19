package util_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	. "util"
)

func TestUnzip(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(tempDir)

	Unzip(filepath.Join("testdata", "test.zip"), tempDir)
	if content, err := ioutil.ReadFile(filepath.Join(tempDir, "test", "textfile.txt")); err == nil {
		println(string(content))
	} else {
		t.Error(err)
	}
}
