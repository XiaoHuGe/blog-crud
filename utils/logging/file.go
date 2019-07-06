package logging

import (
	"fmt"
	"log"
	"os"
	"time"
	"xhblog/utils/setting"
)

func getLogFilePath() string {
	return setting.AppSetting.LogSavePath
}

func getLogFileFullPath() string {
	prefix := getLogFilePath()
	suffix := fmt.Sprintf("%s%s.%s", setting.AppSetting.LogFileName,
		time.Now().Format(setting.AppSetting.TimeFormat), setting.AppSetting.LogFileExt)
	return fmt.Sprintf("%s%s", prefix, suffix)
}

func OpenLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
		case os.IsNotExist(err):
			MkDir()
		case os.IsPermission(err):
			log.Println("os.stat Permission:", err)
	}
	file, err := os.OpenFile(filePath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Println("os.OpenFile err:", err)
	}
	return file
}

func MkDir() {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Println("os.Getwd:", err)
		return
	}
	err = os.MkdirAll(rootDir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		log.Println("os.MkdirAll:", err)
		panic(err)
	}
}
