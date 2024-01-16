package main

import (
	"fmt"
	"os"

	ini "gopkg.in/ini.v1"
)

const CONFIG_PATH = "/.config/local/rfid-reader.ini"

func configGet() (serverPort, readerAddress, readerPort string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Fail to get home directory: %v", err)
		os.Exit(1)
	}

	cfg, err := ini.Load(home + CONFIG_PATH)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// Classic read of values, default section can be represented as empty string
	return cfg.Section("server").Key("port").String(),
		cfg.Section("reader").Key("address").String(),
		cfg.Section("reader").Key("port").String()
}
