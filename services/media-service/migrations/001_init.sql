CREATE TABLE IF NOT EXISTS files (
  id BIGINT NOT NULL AUTO_INCREMENT,
  original_name VARCHAR(255) NOT NULL,
  url VARCHAR(255) NOT NULL,
  path VARCHAR(512) NOT NULL,
  content_type VARCHAR(128) NOT NULL,
  size BIGINT NOT NULL,
  biz_type VARCHAR(64) NOT NULL DEFAULT 'common',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_files_biz_type_created_at (biz_type, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
