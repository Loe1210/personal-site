package domain

import "time"

// User 是认证域对用户信息的最小表示，不暴露密码哈希。
type User struct {
	ID        int64
	Username  string
	Nickname  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
