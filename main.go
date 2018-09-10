package main

import (
	"io"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/remeh/sizedwaitgroup"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//Config for all settings
var Config Configuration

// FolderName separator for files to be parsed
var FolderName string

// SWG Concurency limit
var SWG sizedwaitgroup.SizedWaitGroup

func main() {
	//if len(os.Args) <= 1 {
	//		log.Fatal("Please enter the folderName for files")
	//}
	FolderName = "sample_xlsx" //os.Args[1]

	initLogging()
	parseConfig("config")
	initWorkspace()

	SWG = sizedwaitgroup.New(Config.Concurency.Extract)
	for _, file := range walkPath(Config.Path.Raw + "/" + FolderName) {
		log.Info("Uncompress file: ", file)
		SWG.Add()
		uncompressFile(file)
	}
	SWG.Wait()

	SWG = sizedwaitgroup.New(Config.Concurency.Transform)
	for _, file := range walkPath(Config.Path.Raw + "/" + FolderName) {
		log.Info("Transform file: ", file)
		SWG.Add()
		parseFileConfigs(Config.Path.Config + "/" + FolderName)
		transformFile(file)
	}
	SWG.Wait()

	SWG = sizedwaitgroup.New(Config.Concurency.Loader)
	for _, file := range walkPath(Config.Path.Parsed + "/" + FolderName) {
		log.Info("Load file: ", file)
		SWG.Add()
		parseFileConfigs(Config.Path.Config + "/" + FolderName)
		loadFile(file)
	}
	SWG.Wait()
}

func initWorkspace() {
	os.MkdirAll(Config.Path.Get+"/"+FolderName, os.ModePerm)
	os.MkdirAll(Config.Path.Raw+"/"+FolderName, os.ModePerm)
	os.MkdirAll(Config.Path.Parsed+"/"+FolderName, os.ModePerm)
	os.MkdirAll(Config.Path.Backup+"/"+FolderName, os.ModePerm)
	os.MkdirAll(Config.Path.Config, os.ModePerm)
}

func initLogging() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{})

	logFile, err := os.OpenFile("logs/tigon.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	} else {
		os.MkdirAll("logs", os.ModePerm)
		logFile, err := os.OpenFile("logs/tigon.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err == nil {
			mw := io.MultiWriter(os.Stdout, logFile)
			log.SetOutput(mw)
		} else {
			log.Info("Failed to log to file: ", err)
		}
	}
}

func parseFileConfigs(name string) {
	viper.SetConfigName(name)
	err := viper.ReadInConfig()

	if err != nil {
		log.Warning("Configuration file error " + name + ": " + err.Error())
	}
}

func parseConfig(name string) {
	viper.AddConfigPath(".")
	viper.SetConfigName(name)
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("Configuration file error " + name + ": " + err.Error())
	}

	if viper.Unmarshal(&Config) != nil {
		log.Fatal("Unable to unmarshal config: " + err.Error())
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if viper.Unmarshal(&Config) != nil {
			log.Warning("Failed to unmarshal new config. Using old config: " + err.Error())
		}
	})

}
