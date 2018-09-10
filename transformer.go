package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func transformFile(file string) {
	defer SWG.Done()

	switch getFileExt(file) {
	case "csv":
		err := transformCSV(file)
		if err != nil {
			log.Error("Transform CSV error: ", err.Error())
			return
		}
	case "txt":
		err := transformTXT(file)
		if err != nil {
			log.Error("Transform TXT error: ", err.Error())
			return
		}
	case "xls":
		err := transformXLS(file)
		if err != nil {
			log.Error("Transform XLS error: ", err.Error())
			return
		}
	case "xlsx":
		err := transformXLSX(file)
		if err != nil {
			log.Error("Transform XLSX error: ", err.Error())
			return
		}
	}

	if err := os.Remove(file); err != nil {
		log.Error("Error in deleting the original file: ", err.Error())
	}

}
