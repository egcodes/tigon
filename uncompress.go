package main

import (
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
	log "github.com/sirupsen/logrus"
)

func uncompressFile(file string) {
	defer SWG.Done()

	switch getFileExt(file) {
	case "zip":
		err := archiver.Zip.Open(file, Config.Path.Raw+"/"+FolderName)
		if err != nil {
			log.Error("Uncompress ZIP error: ", err.Error())
			return
		}
	case "tar.gz":
		err := archiver.TarGz.Open(file, Config.Path.Raw+"/"+FolderName)
		if err != nil {
			log.Error("Uncompress TARGZ error: ", err.Error())
			return
		}
	case "tgz":
		err := archiver.TarGz.Open(file, Config.Path.Raw+"/"+FolderName)
		if err != nil {
			log.Error("Uncompress TGZ error: ", err.Error())
			return
		}
	case "tar":
		err := archiver.Tar.Open(file, Config.Path.Raw+"/"+FolderName)
		if err != nil {
			log.Error("Uncompress TAR error: ", err.Error())
			return
		}
	case "tar.bz2":
		err := archiver.TarBz2.Open(file, Config.Path.Raw+"/"+FolderName)
		if err != nil {
			log.Error("Uncompress TARBZ2 error: ", err.Error())
			return
		}
	case "tar.xz":
		err := archiver.TarXZ.Open(file, Config.Path.Raw+"/"+FolderName)
		if err != nil {
			log.Error("Uncompress TARXZ error: ", err.Error())
			return
		}
	case "rar":
		err := archiver.Rar.Open(file, Config.Path.Raw+"/"+FolderName)
		if err != nil {
			log.Error("Uncompress RAR error: ", err.Error())
			return
		}
	case "7z":
		err := uncompress7Z(file)
		if err != nil {
			log.Error("Uncompress 7Z error: ", err.Error())
			return
		}
	case "gz":
		err := uncompressGZ(file)
		if err != nil {
			log.Error("Uncompress GZ error: ", err.Error())
			return
		}
	}

	err := os.Rename(file, Config.Path.Backup+"/"+FolderName+"/"+filepath.Base(file))
	if err != nil {
		log.Error("Error in backup operation: ", err.Error())
	}
}
