package config

import (
	"log"
	"os"
)

var config = map[string]string{
	"SERVICE_NAME":  "log-parser",
	"ENV":           "local",
	"LOG_LEVEL":     "debug",
	"LOG_FILE_PATH": "resources/qgames.log",
}

func GetString(k string) string {
	v := os.Getenv(k)
	if v == "" {
		return config[k]
	}

	return v
}

func LoadFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	return file, nil
}
