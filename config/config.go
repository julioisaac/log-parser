package config

import (
	"log"
	"os"
)

func LoadFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	return file, nil
}
