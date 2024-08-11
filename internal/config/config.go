package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ApplicationInfo struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"applicationInfo"`
	PostgresInfo struct {
		Name        string `json:"name"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		Hostname    string `json:"hostname"`
		Port        int    `json:"port"`
		SSLMode     string `json:"sslmode"`
		SSLRootCert string `json:"sslrootcert"`
		SSLCert     string `json:"sslcert"`
		SSLKey      string `json:"sslkey"`
	} `json:"postgresInfo"`
	AdminCredentials struct {
		Userid   string `json:"userid"`
		Password string `json:"password"`
	} `json:"adminCredentials"`
	RedisConfig struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Password string `json:"password"`
		DbNumber int    `json:"db"`
	} `json:"redisConfig"`
	JwtKey  string `json:"jwtKey"`
	Logging struct {
		Level      string `json:"level"`
		OutputType string `json:"outputType"`
		FilePath   string `json:"filePath"`
	} `json:"logging"`
}

func LoadConfig(configFilePath string) (*Config, error) {
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer configFile.Close()

	var config Config
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	config.overrideWithEnv()

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	return &config, nil
}

func (config *Config) overrideWithEnv() {
	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		config.PostgresInfo.Password = dbPass
	}
	if jwtKey := os.Getenv("JWT_KEY"); jwtKey != "" {
		config.JwtKey = jwtKey
	}
	// Add more environment overrides as needed
}

func (config *Config) validate() error {
	if config.ApplicationInfo.Port == 0 {
		return fmt.Errorf("application port is not set")
	}
	if config.PostgresInfo.Name == "" {
		return fmt.Errorf("database name is not set")
	}
	// Add more validation as needed
	return nil
}

func (config *Config) GetListenAddress() string {
	return fmt.Sprintf("%s:%d", config.ApplicationInfo.Host, config.ApplicationInfo.Port)
}

func (config *Config) GetPostgresConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.PostgresInfo.Hostname,
		config.PostgresInfo.Port,
		config.PostgresInfo.Username,
		config.PostgresInfo.Password,
		config.PostgresInfo.Name,
		config.PostgresInfo.SSLMode)
}
