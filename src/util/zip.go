package util

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(unzipFile string, dest string, options map[string]int) error {
	log.Printf("extract %s to %s", unzipFile, dest)
	rc, err := zip.OpenReader(unzipFile)
	if err != nil {
		return err
	}
	defer rc.Close()

	n := 0
	done := make(chan bool)
	for _, f := range rc.Reader.File {
		// log.Printf("extract file: %s", f.Name)
		if f.FileInfo().IsDir() {
			// log.Printf("extract directory: %s", f.Name)
			// os.MkdirAll(filepath.Join(dest, f.Name), f.Mode())
		} else {
			foc, _ := f.Open()
			defer foc.Close()

			var b bytes.Buffer
			io.Copy(&b, foc)
			// log.Printf("dir is: %s", filepath.Dir(filepath.Join(dest, f.Name)))
			name := ""
			skipComponents := options["skip-components"]
			for _, path := range strings.Split(filepath.Dir(f.Name), string(filepath.Separator))[skipComponents:] {
				name += path + string(filepath.Separator)
			}
			name += filepath.Base(f.Name)

			os.MkdirAll(filepath.Dir(filepath.Join(dest, name)), f.Mode())

			go func(f *zip.File) {
				// log.Printf("extract file: %s", f.Name)
				ioutil.WriteFile(filepath.Join(dest, name), b.Bytes(), f.Mode())
				done <- true
			}(f)
			n++
		}
	}
	for ; n > 0; n-- {
		<-done
	}

	return nil
}
