package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func transformCSV(file string) error {
	fileNameNoExt := getFileNameNoExt(file)

	//Get settings
	parseDataStartIndex := viper.GetInt("transform." + fileNameNoExt + ".parseDataStartIndex")
	parseDataEndIndex := viper.GetInt("transform." + fileNameNoExt + ".parseDataEndIndex")
	parseColumns := viper.GetStringSlice("transform." + fileNameNoExt + ".parseColumns")
	outputSplitChar := viper.GetString("transform." + fileNameNoExt + ".outputSplitChar")

	if parseDataStartIndex == 0 || parseDataEndIndex == 0 || len(parseColumns) == 0 || outputSplitChar == "" {
		log.Warning("No config for this file: " + file)
		return nil
	}

	parsedFile := Config.Path.Parsed + "/" + FolderName + "/" + fileNameNoExt + Config.CustomFileExtention.ParsedFileExt
	f, err := os.Create(parsedFile)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	csvFile, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	rowIndex := 0
	for {
		line, error := reader.Read()

		rowIndex++
		if rowIndex >= parseDataStartIndex {
			if error == io.EOF {
				break
			} else if error != nil {
				return err
			}
			row := ""
			for _, i := range parseColumns {
				columnIndex, err := strconv.Atoi(i)
				if err != nil {
					return err
				}
				if len(line) < columnIndex {
					return err
				}
				row += line[columnIndex] + outputSplitChar
			}

			row = row[:len(row)-1]
			fmt.Fprint(w, row+"\n")
		}

		if rowIndex == parseDataEndIndex {
			break
		}
	}

	csvFile.Close()

	w.Flush()
	f.Sync()
	f.Close()

	return nil
}
