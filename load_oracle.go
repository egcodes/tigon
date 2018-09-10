package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
)

func loadToOracle(file string) error {
	fileName := filepath.Base(file)
	fileNameNoExt := fileName[0 : len(fileName)-len(filepath.Ext(fileName))]

	//Creating control file for sqlldr
	f, err := os.Create(Config.Path.Parsed + "/" + FolderName + "/" + fileNameNoExt + Config.CustomFileExtention.OracleControlFileExt)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	fmt.Fprint(w, viper.GetString("load."+fileNameNoExt+".loadControlFile"))
	w.Flush()
	f.Sync()
	f.Close()

	sqlldrCommand := "sqlldr" + " userid='" + "NORTHI_PARSER_SETTINGS" + "/" + "sarmisak1" + "@" + "NORTHIDB" + "'" +
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
