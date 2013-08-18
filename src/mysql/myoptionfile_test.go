package mysql_test

import (
	"encoding/json"
	. "mysql"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	props, _ := LoadOptions("testdata/my.ini")
	str, _ := json.Marshal(props)
	println(string(str))
	if len(props) == 0 {
		t.FailNow()
	}
}

func TestSave(t *testing.T) {
	props, _ := LoadOptions("testdata/my.ini")
	props["mysqld"]["rpl_semi_sync_master_enabled"] = "true"
	props["mysqlbackup"] = make(map[string]string)
	props["mysqlbackup"]["user"] = "vagrant"
	props["mysqlhotcopy"] = make(map[string]string)

	delete(props["mysql"], "no-auto-rehash")
	delete(props, "client")
	SaveOptions("testdata/my.ini", "testdata/my_new.ini", props)
	defer os.Remove("testdata/my_new.ini")

	props, _ = LoadOptions("testdata/my_new.ini")

	if props["mysqlbackup"]["user"] != "vagrant" {
		t.Fatalf("the user name in the section of mysqlbackup should be %s", "vagrant")
	}

}
