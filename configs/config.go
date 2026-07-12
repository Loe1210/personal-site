package configs

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	MySQL   MySQLConfig   `yaml:"mysql"`
	Session SessionConfig `yaml:"session"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type SessionConfig struct {
	Secret string `yaml:"secret"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

var AppConfig *Config

func Load(configPath string) (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "",
			Port: "8888",
		},
		MySQL: MySQLConfig{
			Host:    "127.0.0.1",
			Port:    "3306",
			User:    "root",
			Charset: "utf8mb4",
		},
		Session: SessionConfig{
			Secret: "personal-site-session-secret",
		},
	}

	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	if data, err := os.ReadFile(configPath); err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// 环境变量覆盖配置文件
	if host := os.Getenv("APP_HOST"); host != "" {
		cfg.Server.Host = host
	}
	if port := os.Getenv("APP_PORT"); port != "" {
		cfg.Server.Port = port
	}
	if host := os.Getenv("MYSQL_HOST"); host != "" {
		cfg.MySQL.Host = host
	}
	if port := os.Getenv("MYSQL_PORT"); port != "" {
		cfg.MySQL.Port = port
	}
	if user := os.Getenv("MYSQL_USER"); user != "" {
		cfg.MySQL.User = user
	}
	if password := os.Getenv("MYSQL_PASSWORD"); password != "" {
		cfg.MySQL.Password = password
	}
	if dbname := os.Getenv("MYSQL_DBNAME"); dbname != "" {
		cfg.MySQL.DBName = dbname
	}
	if charset := os.Getenv("MYSQL_CHARSET"); charset != "" {
		cfg.MySQL.Charset = charset
	}
	if secret := os.Getenv("SESSION_SECRET"); secret != "" {
		cfg.Session.Secret = secret
	}

	AppConfig = cfg
	return cfg, nil
}

func GetServerAddr() string {
	if AppConfig == nil {
		return ":8888"
	}
	return AppConfig.Server.Host + ":" + AppConfig.Server.Port
}
