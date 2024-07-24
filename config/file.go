package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	FileName string
	File     ConfigFile
}

type ConfigFile struct {
	Client ClientConfig `json:"client"`
}

type ClientConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	TenantID     string `json:"tenant_id"`
}

func NewCustomConfig(filename string) *Config {
	return &Config{
		FileName: filename,
	}
}

func (c *Config) Load() (ConfigFile, error) {
	log.Println("Loading config file...")
	file, err := os.Open(c.FileName)
	if err != nil {
		return ConfigFile{}, err
	}
	defer file.Close()
	log.Println("Config file loaded successfully!")

	// Unmarshal the JSON encoding of the file
	log.Println("Decoding config file...")
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c.File)
	if err != nil {
		return ConfigFile{}, err
	}
	log.Println("Config file decoded successfully!")

	return c.File, nil
}

func (c *Config) GetClientConfig() ClientConfig {
	return c.File.Client
}
