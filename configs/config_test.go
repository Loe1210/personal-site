package configs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadUsesDefaultsWhenConfigFileMissing(t *testing.T) {
	t.Setenv("APP_HOST", "")
	t.Setenv("APP_PORT", "")
	t.Setenv("MYSQL_HOST", "")
	t.Setenv("MYSQL_PORT", "")
	t.Setenv("MYSQL_USER", "")
	t.Setenv("MYSQL_PASSWORD", "")
	t.Setenv("MYSQL_DBNAME", "")
	t.Setenv("MYSQL_CHARSET", "")
	t.Setenv("SESSION_SECRET", "")
	t.Setenv("SESSION_STORE_PREFIX", "")
	t.Setenv("SESSION_STORE_EXPIRE_HOUR", "")
	t.Setenv("SESSION_STORE_COOKIE_NAME", "")
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "")
	t.Setenv("UPLOAD_ROOT_DIR", "")
	t.Setenv("UPLOAD_PUBLIC_BASE_PATH", "")
	t.Setenv("UPLOAD_MAX_IMAGE_SIZE_MB", "")
	t.Setenv("SITE_TITLE", "")
	t.Setenv("SITE_BASE_URL", "")

	cfg, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Server.Port != "8888" {
		t.Fatalf("expected default server port 8888, got %q", cfg.Server.Port)
	}
	if cfg.SessionStore.CookieName != "session_id" {
		t.Fatalf("expected default session cookie name session_id, got %q", cfg.SessionStore.CookieName)
	}
	if cfg.SessionStore.Prefix != "session:" {
		t.Fatalf("expected default session prefix session:, got %q", cfg.SessionStore.Prefix)
	}
	if cfg.SessionStore.ExpireHour != 2 {
		t.Fatalf("expected default session expiry 2 hours, got %d", cfg.SessionStore.ExpireHour)
	}
	if cfg.Redis.Addr != "127.0.0.1:6379" {
		t.Fatalf("expected default redis addr, got %q", cfg.Redis.Addr)
	}
	if cfg.Upload.RootDir != "static/uploads/images" {
		t.Fatalf("expected default upload root dir, got %q", cfg.Upload.RootDir)
	}
	if cfg.Upload.PublicBasePath != "/static/uploads/images" {
		t.Fatalf("expected default upload public base path, got %q", cfg.Upload.PublicBasePath)
	}
	if cfg.Upload.MaxImageSizeMB != 5 {
		t.Fatalf("expected default upload max size 5, got %d", cfg.Upload.MaxImageSizeMB)
	}
	if cfg.Site.Title == "" {
		t.Fatal("expected default site title to be set")
	}
}

func TestLoadMergesYamlAndEnvOverrides(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	yaml := []byte("server:\n  host: 127.0.0.1\n  port: \"9999\"\nsession:\n  secret: yaml-secret\nsession_store:\n  prefix: \"yaml-session:\"\n  expire_hour: 6\n  cookie_name: yaml_session\nredis:\n  addr: redis:6379\n  password: redis-pass\n  db: 3\nupload:\n  root_dir: data/uploads\n  public_base_path: /files\n  max_image_size_mb: 9\nsite:\n  title: YAML Title\n  base_url: https://yaml.example.com\nmysql:\n  host: yaml-db\n  port: \"3307\"\n  user: yaml-user\n  password: yaml-pass\n  dbname: yaml-dbname\n  charset: utf8\n")
	if err := os.WriteFile(configPath, yaml, 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	t.Setenv("APP_PORT", "10000")
	t.Setenv("SESSION_STORE_COOKIE_NAME", "env_session")
	t.Setenv("SESSION_STORE_EXPIRE_HOUR", "12")
	t.Setenv("REDIS_ADDR", "env-redis:6379")
	t.Setenv("UPLOAD_ROOT_DIR", "env/uploads")
	t.Setenv("UPLOAD_MAX_IMAGE_SIZE_MB", "12")
	t.Setenv("SITE_BASE_URL", "https://env.example.com")
	t.Setenv("MYSQL_PASSWORD", "env-pass")

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Fatalf("expected YAML server host, got %q", cfg.Server.Host)
	}
	if cfg.Server.Port != "10000" {
		t.Fatalf("expected env override for server port, got %q", cfg.Server.Port)
	}
	if cfg.SessionStore.Prefix != "yaml-session:" {
		t.Fatalf("expected YAML session prefix, got %q", cfg.SessionStore.Prefix)
	}
	if cfg.SessionStore.CookieName != "env_session" {
		t.Fatalf("expected env override for session cookie name, got %q", cfg.SessionStore.CookieName)
	}
	if cfg.SessionStore.ExpireHour != 12 {
		t.Fatalf("expected env override for session expiry, got %d", cfg.SessionStore.ExpireHour)
	}
	if cfg.Redis.Addr != "env-redis:6379" {
		t.Fatalf("expected env override for redis addr, got %q", cfg.Redis.Addr)
	}
	if cfg.Redis.Password != "redis-pass" {
		t.Fatalf("expected YAML redis password, got %q", cfg.Redis.Password)
	}
	if cfg.Redis.DB != 3 {
		t.Fatalf("expected YAML redis db, got %d", cfg.Redis.DB)
	}
	if cfg.Upload.RootDir != "env/uploads" {
		t.Fatalf("expected env override for upload root dir, got %q", cfg.Upload.RootDir)
	}
	if cfg.Upload.PublicBasePath != "/files" {
		t.Fatalf("expected YAML upload public base path, got %q", cfg.Upload.PublicBasePath)
	}
	if cfg.Upload.MaxImageSizeMB != 12 {
		t.Fatalf("expected env override for upload max size, got %d", cfg.Upload.MaxImageSizeMB)
	}
	if cfg.Site.Title != "YAML Title" {
		t.Fatalf("expected YAML site title, got %q", cfg.Site.Title)
	}
	if cfg.Site.BaseURL != "https://env.example.com" {
		t.Fatalf("expected env override for site base url, got %q", cfg.Site.BaseURL)
	}
	if cfg.MySQL.Password != "env-pass" {
		t.Fatalf("expected env override for mysql password, got %q", cfg.MySQL.Password)
	}
}
