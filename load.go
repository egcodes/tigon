package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func loadFile(file string) {
	defer SWG.Done()

	switch viper.GetString("load." + getFileNameNoExt(file) + ".db") {
	case "oracle":
		err := loadToOracle(file)
		if err != nil {
			log.Error("Load Oracle error: ", err.Error())
			return
		}
	}

	if err := os.Remove(file); err != nil {
		log.Error("Errors in deleting the parsed file: ", err.Error())
	}
	if err := os.Remove(strings.Replace(file, Config.CustomFileExtention.ParsedFileExt, Config.CustomFileExtention.OracleControlFileExt, -1)); err != nil {
		log.Error("Errors in deleting the parsed file: ", err.Error())
	}

}
