package configs

import (
	"github.com/Loe1210/personal-site/pkg/xconfig"
)

type Config struct {
	Server       ServerConfig       `yaml:"server"`
	MySQL        MySQLConfig        `yaml:"mysql"`
	Session      SessionConfig      `yaml:"session"`
	SessionStore SessionStoreConfig `yaml:"session_store"`
	Redis        RedisConfig        `yaml:"redis"`
	Upload       UploadConfig       `yaml:"upload"`
	Site         SiteConfig         `yaml:"site"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type SessionConfig struct {
	Secret string `yaml:"secret"`
}

type SessionStoreConfig struct {
	Prefix     string `yaml:"prefix"`
	ExpireHour int    `yaml:"expire_hour"`
	CookieName string `yaml:"cookie_name"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
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
	cfg := defaultConfig()

	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	if err := xconfig.LoadYAML(configPath, cfg); err != nil {
		return nil, err
	}

	xconfig.OverrideString(&cfg.Server.Host, "APP_HOST")
	xconfig.OverrideString(&cfg.Server.Port, "APP_PORT")
	xconfig.OverrideString(&cfg.MySQL.Host, "MYSQL_HOST")
	xconfig.OverrideString(&cfg.MySQL.Port, "MYSQL_PORT")
	xconfig.OverrideString(&cfg.MySQL.User, "MYSQL_USER")
	xconfig.OverrideString(&cfg.MySQL.Password, "MYSQL_PASSWORD")
	xconfig.OverrideString(&cfg.MySQL.DBName, "MYSQL_DBNAME")
	xconfig.OverrideString(&cfg.MySQL.Charset, "MYSQL_CHARSET")
	xconfig.OverrideString(&cfg.Session.Secret, "SESSION_SECRET")
	xconfig.OverrideString(&cfg.SessionStore.Prefix, "SESSION_STORE_PREFIX")
	xconfig.OverrideInt(&cfg.SessionStore.ExpireHour, "SESSION_STORE_EXPIRE_HOUR")
	xconfig.OverrideString(&cfg.SessionStore.CookieName, "SESSION_STORE_COOKIE_NAME")
	xconfig.OverrideString(&cfg.Redis.Addr, "REDIS_ADDR")
	xconfig.OverrideString(&cfg.Redis.Password, "REDIS_PASSWORD")
	xconfig.OverrideInt(&cfg.Redis.DB, "REDIS_DB")
	xconfig.OverrideString(&cfg.Upload.RootDir, "UPLOAD_ROOT_DIR")
	xconfig.OverrideString(&cfg.Upload.PublicBasePath, "UPLOAD_PUBLIC_BASE_PATH")
	xconfig.OverrideInt64(&cfg.Upload.MaxImageSizeMB, "UPLOAD_MAX_IMAGE_SIZE_MB")
	xconfig.OverrideString(&cfg.Site.Title, "SITE_TITLE")
	xconfig.OverrideString(&cfg.Site.BaseURL, "SITE_BASE_URL")

	if cfg.SessionStore.Prefix == "" {
		cfg.SessionStore.Prefix = "session:"
	}
	if cfg.SessionStore.CookieName == "" {
		cfg.SessionStore.CookieName = "session_id"
	}
	if cfg.SessionStore.ExpireHour <= 0 {
		cfg.SessionStore.ExpireHour = 2
	}
	if cfg.Upload.MaxImageSizeMB <= 0 {
		cfg.Upload.MaxImageSizeMB = 5
	}

	AppConfig = cfg
	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
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
		SessionStore: SessionStoreConfig{
			Prefix:     "session:",
			ExpireHour: 2,
			CookieName: "session_id",
		},
		Redis: RedisConfig{
			Addr: "127.0.0.1:6379",
			DB:   0,
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
}

func GetServerAddr() string {
	if AppConfig == nil {
		return ":8888"
	}
	return AppConfig.Server.Host + ":" + AppConfig.Server.Port
}
