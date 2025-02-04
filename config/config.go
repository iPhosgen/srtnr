package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Service  ServiceConfig  `yaml:"service"`
}

type ServiceConfig struct {
	Host string `yaml:"host" env:"SERVICE_HOST,required"`
	Port int    `yaml:"port" env:"SERVICE_PORT,required"`
}

type DatabaseConfig struct {
	Host     string             `yaml:"host" env:"DATABASE_HOST,required"`
	Port     int                `yaml:"port" env:"DATABASE_PORT,required"`
	User     string             `yaml:"user" env:"DATABASE_USR,required"`
	Password string             `yaml:"password" env:"DATABASE_PWD,required"`
	DBName   string             `yaml:"dbname" env:"DATABASE_DBNAME,required"`
	SSL      *DatabaseSSLConfig `yaml:"ssl"`
}

type DatabaseSSLConfig struct {
	SSLMode  bool   `yaml:"sslmode" env:"DATABASE_SSL_MODE"`
	CertFile string `yaml:"cert_file" env:"DATABASE_SSL_CERT"`
	KeyFile  string `yaml:"key_file" env:"DATABASE_SSL_KEY"`
}

func (dc *DatabaseConfig) BuildDSN() (dsn string) {
	dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		dc.User, dc.Password, dc.Host, dc.Port, dc.DBName)

	if dc.SSL != nil {
		if dc.SSL.SSLMode {
			sslOptions := []string{}
			if len(dc.SSL.CertFile) > 0 {
				sslOptions = append(sslOptions, fmt.Sprintf("sslcert=%s", dc.SSL.CertFile))
			}
			if len(dc.SSL.KeyFile) > 0 {
				sslOptions = append(sslOptions, fmt.Sprintf("sslkey=%s", dc.SSL.KeyFile))
			}

			dsn = fmt.Sprintf("%s?%s", dsn, strings.Join(sslOptions, "&"))
		}
	}

	return
}

func LoadConfig(filePath string) (*Config, error) {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Printf("failed to read config from ENV: %v", err)
	} else {
		return cfg, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("invalid config file: %v", err)
		return nil, err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		log.Fatalf("config unmarshal error: %v", err)
		return nil, err
	}

	if cfg.Database.SSL != nil && cfg.Database.SSL.SSLMode {
		if _, err := os.Stat(cfg.Database.SSL.CertFile); err != nil {
			log.Fatalf("cert file not found: %v", err)
			return nil, err
		}

		if _, err := os.Stat(cfg.Database.SSL.KeyFile); err != nil {
			log.Fatalf("key file not found: %v", err)
			return nil, err
		}
	}

	return cfg, nil
}
