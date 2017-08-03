package config

import (
	"os"
	"log"
	"encoding/json"
)

var D Data

type Data struct {
	Server	 			string	`json:"server"`
	DashboardTemplate 	string	`json:"dashboard_template"`
	ContentDir 			string	`json:"content_dir"`
	StoragePath 		string	`json:"storage_path"`
}

func LoadConfig(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Config file is missing: ", path)
	}

	if err := json.NewDecoder(f).Decode(&D); err != nil {
		log.Fatal(err)
	}
}