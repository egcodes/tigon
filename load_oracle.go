package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func loadToOracle(file string) error {
	fileName := filepath.Base(file)
	fileNameNoExt := fileName[0 : len(fileName)-len(filepath.Ext(fileName))]

	username := viper.GetString("load." + fileNameNoExt + ".userName")
	password := viper.GetString("load." + fileNameNoExt + ".password")
	tnsname := viper.GetString("load." + fileNameNoExt + ".tnsName")
	ctlFileContent := viper.GetString("load." + fileNameNoExt + ".loadControlFile")

	if username == "" || password == "" || tnsname == "" || ctlFileContent == "" {
		log.Warning("Load: No config for this file: " + file)
		return nil
	}

	//Creating control file for sqlldr
	f, err := os.Create(Config.Path.Parsed + "/" + FolderName + "/" + fileNameNoExt + Config.CustomFileExtention.OracleControlFileExt)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	fmt.Fprint(w, ctlFileContent)
	w.Flush()
	f.Sync()
	f.Close()

	sqlldrCommand := "sqlldr" + " userid='" + username + "/" + password + "@" + tnsname + "'" +
		" control='" + Config.Path.Parsed + "/" + FolderName + "/" + fileNameNoExt + ".ctl" + "'" +
		" log='" + Config.Path.Parsed + "/" + FolderName + "/" + fileNameNoExt + ".log" + "'" +
		" direct=true rows=" + "100000" +
		" errors=" + "100"

	cmd := exec.Command("sh", "-c", sqlldrCommand)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
