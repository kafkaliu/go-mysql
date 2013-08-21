package mysql

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"util"
)

const (
	mysqlservice = "mysql-ha"
)

func Extract(installPackage string, optionFile string) error {
	if props, err := LoadOptions(optionFile); err == nil {
		switch {
		case strings.HasSuffix(installPackage, ".zip"):
			if _, err = os.Stat(props["mysqld"]["basedir"]); err == nil {
				return errors.New("Fail to extract due to base dir exists.")
			}
			if _, err = os.Stat(props["mysqld"]["datadir"]); err == nil {
				return errors.New("Fail to extract due to data dir exists.")
			}
			basedir := props["mysqld"]["basedir"]
			log.Printf("install package from %s to %s", installPackage, basedir)
			if err = os.Mkdir(basedir, 0644); err != nil {
				return err
			}
			options := map[string]int{"skip-components": 1}
			if err = util.Unzip(installPackage, basedir, options); err != nil {
				return err
			}
			b, _ := ioutil.ReadFile(optionFile)
			os.MkdirAll(filepath.Join(basedir, "etc"), 0644)
			ioutil.WriteFile(filepath.Join(basedir, "etc", "my.ini"), b, 0644)
		case strings.HasSuffix(installPackage, ".tar.gz"):

		}
		return nil
	} else {
		return err
	}
}

func mysql(basedir, sql string) error {
	return execute(filepath.Join(basedir, "bin", "mysql"), "-uroot", "-h127.0.0.1", "-e", sql)
}

func execute(name string, arg ...string) error {
	out, err := exec.Command(name, arg...).Output()
	log.Printf("%s", out)
	log.Printf("%v", err)
	return err
}

func InstallServer(basedir string, datadir string, dataurl string, defaultdata string) error {
	if _, err := os.Stat(datadir); err == nil {
		return errors.New("Fail to extract due to data dir exists.")
	}
	os.MkdirAll(datadir, 0644)
	if dataurl == "" {
		options := map[string]int{"skip-components": 1}
		if err := util.Unzip(defaultdata, datadir, options); err != nil {
			return err
		}
	} else {
		// todo: import data from dataurl
	}

	switch runtime.GOOS {
	case "windows":
		defaultsFile := filepath.Join(basedir, "etc", "my.ini")
		if err := execute(filepath.Join(basedir, "bin", "mysqld"), "--install", "mysql-ha", "--defaults-file="+defaultsFile); err != nil {
			return err
		}

		rpl_semi_sync_master_enabled := "0"
		rpl_semi_sync_slave_enabled := "0"
		if props, err := LoadOptions(defaultsFile); err == nil {
			if props["mysqld"]["rpl_semi_sync_master_enabled"] == "1" {
				rpl_semi_sync_master_enabled = "1"
				delete(props["mysqld"], "rpl_semi_sync_master_enabled")
			}
			if props["mysqld"]["rpl_semi_sync_slave_enabled"] == "1" {
				rpl_semi_sync_slave_enabled = "1"
				delete(props["mysqld"], "rpl_semi_sync_slave_enabled")
			}
			if err = SaveOptions(defaultsFile, defaultsFile, props); err != nil {
				return err
			}
		} else {
			return err
		}
		// log.Printf("rpl_semi_sync_master_enabled = %s, rpl_semi_sync_slave_enabled = %s", rpl_semi_sync_master_enabled, rpl_semi_sync_slave_enabled)

		if err := Start(basedir); err != nil {
			return err
		}
		if err := mysql(basedir, "delete from mysql.user where user=''"); err != nil {
			return err
		}
		if err := mysql(basedir, "grant replication slave, replication client on *.* to repl identified by 'repl'"); err != nil {
			return err
		}
		if err := mysql(basedir, "grant replication slave, replication client on *.* to 'repl'@'127.0.0.1' identified by 'repl'"); err != nil {
			return err
		}
		if err := mysql(basedir, "install plugin rpl_semi_sync_master soname 'semisync_master.dll'"); err != nil {
			return err
		}
		if err := mysql(basedir, "install plugin rpl_semi_sync_slave soname 'semisync_slave.dll'"); err != nil {
			return err
		}
		if err := mysql(basedir, "set global rpl_semi_sync_master_enabled = 1"); err != nil {
			return err
		}
		if err := mysql(basedir, "set global rpl_semi_sync_slave_enabled = 1"); err != nil {
			return err
		}

		if props, err := LoadOptions(defaultsFile); err == nil {
			props["mysqld"]["rpl_semi_sync_master_enabled"] = rpl_semi_sync_master_enabled
			props["mysqld"]["rpl_semi_sync_slave_enabled"] = rpl_semi_sync_slave_enabled
			if err = SaveOptions(defaultsFile, defaultsFile, props); err != nil {
				return err
			}
		} else {
			return err
		}

	default:
	}
	// switch runtime.GOOS == "windows" {

	// }
	return nil
}

func UninstallServer(basedir string, datadir string) error {
	Stop(basedir)
	defer os.RemoveAll(basedir)
	defer os.RemoveAll(datadir)
	switch runtime.GOOS {
	case "windows":
		if err := execute("sc", "delete", mysqlservice); err != nil {
			return err
		}
	default:
	}
	return nil
}

func Start(basedir string) error {
	switch runtime.GOOS {
	case "windows":
		if err := execute("net", "start", mysqlservice); err != nil {
			return err
		}
	}
	return nil
}

func Stop(basedir string) error {
	switch runtime.GOOS {
	case "windows":
		if err := execute("net", "stop", mysqlservice); err != nil {
			return err
		}
	}
	return nil
}
