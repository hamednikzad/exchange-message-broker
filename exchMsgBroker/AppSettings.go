package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//AppSettings values
type AppSettings struct {
	Logging Logging `json:"logging"`
	Server  Server  `json:"server"`
}

//Logging settings
type Logging struct {
	Level string `json:"level"`
}
type Server struct {
	Listen string `json:"listen"`
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

const appSettingName = "./appsettings.json"

//GetAppSettings settings
func getAppSettings() AppSettings {
	currWd, _ := os.Getwd()
	settingsPath := filepath.Join(currWd, appSettingName)
	if !exists(settingsPath) {
		settingsPath = filepath.Join(currWd, "exchMsgBroker", appSettingName)
	}
	jsonFile, err := os.Open(settingsPath)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var appSettings AppSettings
	json.Unmarshal(byteValue, &appSettings)
	fmt.Println("Appsettings: ", appSettings)
	return appSettings
}
