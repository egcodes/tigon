package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func copyFile(file string) error {
	from, err := os.Open(file)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(Config.Path.Raw+"/"+FolderName+"/"+filepath.Base(file), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}

	return nil
}

func walkPath(path string) []string {
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != Config.CustomFileExtention.OracleControlFileExt && !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return files
}

func getFileExt(file string) string {
	return strings.SplitAfterN(file, ".", 2)[1]
}

func getFileNameNoExt(file string) string {
	fileName := filepath.Base(file)
	fileNameNoExt := fileName[0 : len(fileName)-len(filepath.Ext(fileName))]
	return fileNameNoExt
}
