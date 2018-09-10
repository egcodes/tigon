package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tealeg/xlsx"
)

func transformXLSX(file string) error {
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

	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		return err
	}

	sheet := xlFile.Sheets[parseSheet]

	rowIndex := 0
	for _, line := range sheet.Rows {
		rowIndex++

		if rowIndex >= parseDataStartIndex {
			row := ""
			for _, i := range parseColumns {
				columnIndex, err := strconv.Atoi(i)
				if err != nil {
					return err
				}

				cells := line.Cells
				if len(cells) > columnIndex {
					cell := cells[columnIndex]
					var s interface{}

					switch cell.Type() {
					case xlsx.CellTypeNumeric:
						s, _ = cell.Float()
					case xlsx.CellTypeString:
						s = cell.String()
					case xlsx.CellTypeFormula:
						s = cell.Formula()
					case xlsx.CellTypeBool:
						s = cell.Bool()
					default:
						s = ""
					}
					row += fmt.Sprintf("%v%s", s, outputSplitChar)
				} else {
					row += "" + outputSplitChar
				}
			}
			row = row[:len(row)-1]
			fmt.Fprint(w, row+"\n")
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
