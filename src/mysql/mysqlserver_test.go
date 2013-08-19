package mysql_test

import (
	"io/ioutil"
	"log"
	. "mysql"
	"os"
	"testing"
)

func TestUncompressAndRemoveData(t *testing.T) {
	defer os.RemoveAll("c:/temp/mysql-install")
	defer os.RemoveAll("c:/temp/mysql-data")
	if err := UncompressAndRemoveData("testdata/mysql-5.5.33-winx64.zip", "testdata/my.ini"); err != nil {
		t.Fatalf("Fail to uncompress mysql.", err)
	}
	if _, err := os.Stat("c:/temp/mysql-install"); err != nil {
		t.Fatalf("Fail to uncompress.")
	}
	if _, err := os.Stat("c:/temp/mysql-data"); err == nil {
		t.Fatalf("Should remove data dir.")
	}
	if content, err := ioutil.ReadFile("c:/temp/mysql-install/mysql-5.5.33-winx64/README"); err == nil {
		log.Println(string(content))
	} else {
		t.Error(err)
	}
}
