package mysql

import (
	"log"
	"os"
	"strings"
	"util"
)

func UncompressAndRemoveData(installPackage string, optionFile string) error {
	if props, err := LoadOptions(optionFile); err == nil {
		switch {
		case strings.HasSuffix(installPackage, ".zip"):
			basedir := props["mysqld"]["basedir"]
			log.Printf("install package from %s to %s", installPackage, basedir)
			if err = os.Mkdir(basedir, 0644); err != nil {
				return err
			}
			if err = util.Unzip(installPackage, basedir); err != nil {
				return err
			}
			os.RemoveAll(props["mysqld"]["datadir"])
		case strings.HasSuffix(installPackage, ".tar.gz"):

		}
		return nil
	} else {
		return err
	}
}
