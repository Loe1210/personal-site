package configs

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	MySQL   MySQLConfig   `yaml:"mysql"`
	Session SessionConfig `yaml:"session"`
	Upload  UploadConfig  `yaml:"upload"`
	Site    SiteConfig    `yaml:"site"`
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

type UploadConfig struct {
	RootDir        string `yaml:"root_dir"`
	PublicBasePath string `yaml:"public_base_path"`
	MaxImageSizeMB int64  `yaml:"max_image_size_mb"`
}

type SiteConfig struct {
	Title   string `yaml:"title"`
	BaseURL string `yaml:"base_url"`
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
		Upload: UploadConfig{
			RootDir:        "static/uploads/images",
			PublicBasePath: "/static/uploads/images",
			MaxImageSizeMB: 5,
		},
		Site: SiteConfig{
			Title:   "Personal Site",
			BaseURL: "http://localhost:8888",
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
	if rootDir := os.Getenv("UPLOAD_ROOT_DIR"); rootDir != "" {
		cfg.Upload.RootDir = rootDir
	}
	if publicBasePath := os.Getenv("UPLOAD_PUBLIC_BASE_PATH"); publicBasePath != "" {
		cfg.Upload.PublicBasePath = publicBasePath
	}
	if maxImageSizeMB := os.Getenv("UPLOAD_MAX_IMAGE_SIZE_MB"); maxImageSizeMB != "" {
		if parsed, err := strconv.ParseInt(maxImageSizeMB, 10, 64); err == nil && parsed > 0 {
			cfg.Upload.MaxImageSizeMB = parsed
		}
	}
	if title := os.Getenv("SITE_TITLE"); title != "" {
		cfg.Site.Title = title
	}
	if baseURL := os.Getenv("SITE_BASE_URL"); baseURL != "" {
		cfg.Site.BaseURL = baseURL
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
