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
		cmd := exec.Command(filepath.Join(basedir, "bin", "mysqld"), "--install", "mysql-ha", "--defaults-file="+filepath.Join(basedir, "etc", "my.ini"))
		out, err := cmd.Output()
		log.Println(string(out))
		if err != nil {
			return err
		}
		return Start(basedir)
	default:
	}
	// switch runtime.GOOS == "windows" {

	// }
	return nil
}

func UninstallServer(basedir string, datadir string) error {
	Stop(basedir)
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command(filepath.Join(basedir, "bin", "mysqld"), "--remove", mysqlservice)
		out, err := cmd.Output()
		log.Println(string(out))
		if err != nil {
			return err
		}
	default:
	}

	os.RemoveAll(basedir)
	os.RemoveAll(datadir)
	return nil
}

func Start(basedir string) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("net", "start", mysqlservice)
		out, err := cmd.Output()
		log.Println(string(out))
		if err != nil {
			return err
		}
	}
	return nil
}

func Stop(basedir string) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("net", "stop", mysqlservice)
		out, err := cmd.Output()
		log.Println(string(out))
		if err != nil {
			return err
		}
	}
	return nil
}
