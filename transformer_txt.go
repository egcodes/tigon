package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func transformTXT(file string) error {
	fileNameNoExt := getFileNameNoExt(file)

	//Get settings
	parseDataStartIndex := viper.GetInt("transform." + fileNameNoExt + ".parseDataStartIndex")
	parseDataEndIndex := viper.GetInt("transform." + fileNameNoExt + ".parseDataEndIndex")
	parseColumns := viper.GetStringSlice("transform." + fileNameNoExt + ".parseColumns")
	fileSplitChar := viper.GetString("transform." + fileNameNoExt + ".fileSplitChar")
	fileRegexStr := viper.GetString("transform." + fileNameNoExt + ".fileRegexStr")
	outputSplitChar := viper.GetString("transform." + fileNameNoExt + ".outputSplitChar")

	if parseDataStartIndex == 0 || parseDataEndIndex == 0 || len(parseColumns) == 0 ||
		(fileSplitChar == "" && fileRegexStr == "") || outputSplitChar == "" {
		log.Warning("No config for this file: " + file)
		return nil
	}

	parsedFile := Config.Path.Parsed + "/" + FolderName + "/" + fileNameNoExt + Config.CustomFileExtention.ParsedFileExt
	f, err := os.Create(parsedFile)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	inFile, _ := os.Open(file)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	regex := *regexp.MustCompile(fileRegexStr)
	rowIndex := 0
	for scanner.Scan() {
		line := scanner.Text()
		rowIndex++

		if rowIndex >= parseDataStartIndex {

			if fileSplitChar != "" {
				row := ""
				splittedRow := strings.Split(line, fileSplitChar)
				for _, i := range parseColumns {
					columnIndex, err := strconv.Atoi(i)
					if err != nil {
						return err
					}
					if len(line) < columnIndex {
						return err
					}
					row += splittedRow[columnIndex] + outputSplitChar
				}
				row = row[:len(row)-1]
				fmt.Fprint(w, row+"\n")
			} else {
				row := ""
				splittedRow := regex.FindAllStringSubmatch(line, -1)
				for _, i := range parseColumns {
					columnIndex, err := strconv.Atoi(i)
					if err != nil {
						return err
					}
					if len(row) < columnIndex {
						return err
					}
					row += splittedRow[columnIndex][0] + outputSplitChar
				}
				row = row[:len(row)-1]
				fmt.Fprint(w, row+"\n")
			}

		}

		if rowIndex == parseDataEndIndex {
			break
		}
	}
	w.Flush()
	f.Sync()
	f.Close()

	return nil
}
