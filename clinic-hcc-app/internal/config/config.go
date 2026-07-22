package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ClinicNameEn    string `json:"clinic_name_en"`
	ClinicNameZh    string `json:"clinic_name_zh"`
	ClinicAddress   string `json:"clinic_address"`
	ClinicTelephone string `json:"clinic_telephone"`
	ServerPort      int    `json:"server_port"`
	BindAddress     string `json:"bind_address"`
	DatabasePath    string `json:"database_path"`
	SessionSecret   string `json:"session_secret"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
