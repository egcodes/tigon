package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/extrame/xls"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func transformXLS(file string) error {
	fileNameNoExt := getFileNameNoExt(file)

	//Get settings
	parseSheet := viper.GetInt("transform." + fileNameNoExt + ".parseSheet")
	parseDataStartIndex := viper.GetInt("transform." + fileNameNoExt + ".parseDataStartIndex")
	parseDataEndIndex := viper.GetInt("transform." + fileNameNoExt + ".parseDataEndIndex")
	parseColumns := viper.GetStringSlice("transform." + fileNameNoExt + ".parseColumns")
	outputSplitChar := viper.GetString("transform." + fileNameNoExt + ".outputSplitChar")

	if viper.IsSet("transform."+fileNameNoExt+".parseSheet") == false || parseDataStartIndex == 0 || parseDataEndIndex == 0 || len(parseColumns) == 0 || outputSplitChar == "" {
		log.Warning("No config for this file: " + file)
		return nil
	}

	parsedFile := Config.Path.Parsed + "/" + FolderName + "/" + fileNameNoExt + Config.CustomFileExtention.ParsedFileExt
	f, err := os.Create(parsedFile)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	if xlFile, err := xls.Open(file, "utf-8"); err == nil {
		if sheet := xlFile.GetSheet(parseSheet); sheet != nil {
			rowIndex := 0
			for i := 0; i <= (int(sheet.MaxRow)); i++ {
				line := sheet.Row(i)
				rowIndex++

				if rowIndex >= parseDataStartIndex {
					row := ""
					for _, i := range parseColumns {
						columnIndex, err := strconv.Atoi(i)
						if err != nil {
							return err
						}
						if line.LastCol() < columnIndex {
							return err
						}
						row += line.Col(columnIndex) + outputSplitChar
					}

					row = row[:len(row)-1]
					fmt.Fprint(w, row+"\n")
				}

				if rowIndex == parseDataEndIndex {
					break
				}
			}
		} else {
			return err
		}
	} else {
		return err
	}

	w.Flush()
	f.Sync()
	f.Close()

	return nil
}
