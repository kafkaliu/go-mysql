package util

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Unzip(unzipFile string, dest string) error {
	log.Printf("extract %s to %s", unzipFile, dest)
	rc, err := zip.OpenReader(unzipFile)
	if err != nil {
		return err
	}
	defer rc.Close()

	n := 0
	done := make(chan bool)
	for _, f := range rc.Reader.File {
		foc, _ := f.Open()
		defer foc.Close()

		var b bytes.Buffer
		io.Copy(&b, foc)
		if f.FileInfo().IsDir() {
			log.Printf("extract directory: %s", f.Name)
			os.MkdirAll(filepath.Join(dest, f.Name), f.Mode())
		} else {
			log.Printf("extract file: %s", f.Name)
			// log.Printf("dir is: %s", filepath.Dir(filepath.Join(dest, f.Name)))
			os.MkdirAll(filepath.Dir(filepath.Join(dest, f.Name)), f.Mode())

			go func() {
				ioutil.WriteFile(filepath.Join(dest, f.Name), b.Bytes(), f.Mode())
				done <- true
				n++
			}()
			if err = ioutil.WriteFile(filepath.Join(dest, f.Name), b.Bytes(), f.Mode()); err != nil {
				return err
			}
		}
	}
	for ; n > 0; n-- {
		<-done
	}

	return nil
}
