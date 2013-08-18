package mysql

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	my := new(MySqlOptionFile)
	my.FileName = "testdata/my.ini"
	my.Load()
	props := my.Properties
	println(my.String())
	if len(props) == 0 {
		t.FailNow()
	}
}

func TestSave(t *testing.T) {
	my := new(MySqlOptionFile)
	my.FileName = "testdata/my.ini"
	my.Load()
	my.Properties["mysqld"]["rpl_semi_sync_master_enabled"] = "true"
	my.Properties["mysqlbackup"] = make(map[string]string)
	my.Properties["mysqlbackup"]["user"] = "vagrant"
	my.Properties["mysqlhotcopy"] = make(map[string]string)

	delete(my.Properties["mysql"], "no-auto-rehash")
	delete(my.Properties, "client")
	my.Save("testdata/my_new.ini")

	newMy := new(MySqlOptionFile)
	newMy.FileName = "testdata/my_new.ini"
	newMy.Load()

	if newMy.Properties["mysqlbackup"]["user"] != "vagrant" {
		t.Fatalf("the user name in the section of mysqlbackup should be %s", "vagrant")
	}

	defer os.Remove("testdata/my_new.ini")
}
