package mysql_test

import (
	"io/ioutil"
	"log"
	. "mysql"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestExtract(t *testing.T) {
	defer os.RemoveAll("c:/temp/mysql-install")
	defer os.RemoveAll("c:/temp/mysql-data")
	if err := Extract("testdata/mysql-5.5.33-winx64.zip", "testdata/my.ini"); err != nil {
		t.Fatalf("Fail to uncompress mysql.", err)
	}
	if _, err := os.Stat("c:/temp/mysql-install"); err != nil {
		t.Fatalf("Fail to uncompress.")
	}
	if _, err := os.Stat("c:/temp/mysql-data"); err == nil {
		t.Fatalf("Should remove data dir.")
	}
	if content, err := ioutil.ReadFile("c:/temp/mysql-install/etc/my.ini"); err == nil {
		log.Println(string(content))
	} else {
		t.Error(err)
	}
}

func TestExtractInstallUninstallServer(t *testing.T) {
	defer UninstallServer("c:/temp/mysql-install", "c:/temp/mysql-data")
	if err := Extract("testdata/mysql-5.5.33-winx64.zip", "testdata/my.ini"); err != nil {
		t.Fatalf("Fail to uncompress mysql.", err)
	}
	InstallServer("c:/temp/mysql-install", "c:/temp/mysql-data", "", "testdata/mysql-5.5.33-data-winx64.zip")
	if content, err := ioutil.ReadFile("c:/temp/mysql-install/etc/my.ini"); err == nil {
		log.Println(string(content))
	} else {
		t.Error(err)
	}
	if out, err := exec.Command(filepath.Join("c:/temp/mysql-install", "bin", "mysql"), "-uroot", "-h127.0.0.1", "-e", "select 1").Output(); err != nil {
		log.Println(out)
		t.Error(err)
	}
}
