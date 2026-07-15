CREATE TABLE IF NOT EXISTS upload_tasks (
  upload_id VARCHAR(64) NOT NULL,
  user_id BIGINT NOT NULL,
  biz_type VARCHAR(64) NOT NULL DEFAULT 'common',
  biz_id VARCHAR(128) NOT NULL DEFAULT '',
  file_name VARCHAR(255) NOT NULL,
  file_size BIGINT NOT NULL,
  chunk_size BIGINT NOT NULL,
  chunk_count INT NOT NULL,
  uploaded_chunks VARCHAR(4096) NOT NULL DEFAULT '',
  status VARCHAR(32) NOT NULL,
  sha256 VARCHAR(64) NOT NULL DEFAULT '',
  expires_at DATETIME NOT NULL,
  last_error TEXT NULL,
  version BIGINT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (upload_id),
  KEY idx_upload_tasks_user_status (user_id, status),
  KEY idx_upload_tasks_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS upload_chunks (
  upload_id VARCHAR(64) NOT NULL,
  chunk_index INT NOT NULL,
  size BIGINT NOT NULL,
  sha256 VARCHAR(64) NOT NULL DEFAULT '',
  storage_path VARCHAR(512) NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (upload_id, chunk_index),
  CONSTRAINT fk_upload_chunks_task FOREIGN KEY (upload_id) REFERENCES upload_tasks(upload_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE files
  ADD COLUMN upload_id VARCHAR(64) NULL,
  ADD COLUMN sha256 VARCHAR(64) NOT NULL DEFAULT '',
  ADD COLUMN biz_id VARCHAR(128) NOT NULL DEFAULT '',
  ADD KEY idx_files_upload_id (upload_id);
