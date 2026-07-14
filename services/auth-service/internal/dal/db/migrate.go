package db

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"type:varchar(64);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Nickname     string `gorm:"type:varchar(64);not null"`
	Type         string `gorm:"type:varchar(32);index;default:'admin'"`
	Status       string `gorm:"type:varchar(32);index;default:'active'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Role struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(64);not null"`
	Code        string `gorm:"type:varchar(64);uniqueIndex;not null"`
	Description string `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Permission struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"type:varchar(64);not null"`
	Code        string `gorm:"type:varchar(64);uniqueIndex;not null"`
	Resource    string `gorm:"type:varchar(64);index;not null"`
	Action      string `gorm:"type:varchar(64);index;not null"`
	Description string `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserRole struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	UserID    int64 `gorm:"uniqueIndex:idx_user_role;not null"`
	RoleID    int64 `gorm:"uniqueIndex:idx_user_role;not null"`
	CreatedAt time.Time
}

type RolePermission struct {
	ID           int64 `gorm:"primaryKey;autoIncrement"`
	RoleID       int64 `gorm:"uniqueIndex:idx_role_permission;not null"`
	PermissionID int64 `gorm:"uniqueIndex:idx_role_permission;not null"`
	CreatedAt    time.Time
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Role{}, &Permission{}, &UserRole{}, &RolePermission{})
}

func Seed(db *gorm.DB) error {
	superAdminRole, err := ensureRole(db, "Super Admin", "super_admin", "super administrator")
	if err != nil {
		return err
	}

	editorRole, err := ensureRole(db, "Editor", "editor", "content editor")
	if err != nil {
		return err
	}

	permissions := []Permission{
		{Name: "Article Read", Code: "article:read", Resource: "article", Action: "read", Description: "read article"},
		{Name: "Article Create", Code: "article:create", Resource: "article", Action: "create", Description: "create article"},
		{Name: "Article Update", Code: "article:update", Resource: "article", Action: "update", Description: "update article"},
		{Name: "Article Delete", Code: "article:delete", Resource: "article", Action: "delete", Description: "delete article"},
		{Name: "Category Read", Code: "category:read", Resource: "category", Action: "read", Description: "read category"},
		{Name: "Category Create", Code: "category:create", Resource: "category", Action: "create", Description: "create category"},
		{Name: "Category Update", Code: "category:update", Resource: "category", Action: "update", Description: "update category"},
		{Name: "Category Delete", Code: "category:delete", Resource: "category", Action: "delete", Description: "delete category"},
		{Name: "Tag Read", Code: "tag:read", Resource: "tag", Action: "read", Description: "read tag"},
		{Name: "Tag Create", Code: "tag:create", Resource: "tag", Action: "create", Description: "create tag"},
		{Name: "Tag Update", Code: "tag:update", Resource: "tag", Action: "update", Description: "update tag"},
		{Name: "Tag Delete", Code: "tag:delete", Resource: "tag", Action: "delete", Description: "delete tag"},
		{Name: "User Me", Code: "user:me", Resource: "user", Action: "me", Description: "get current user"},
		{Name: "User Logout", Code: "user:logout", Resource: "user", Action: "logout", Description: "logout user"},
	}

	permissionMap := make(map[string]Permission)
	for _, item := range permissions {
		perm, err := ensurePermission(db, item)
		if err != nil {
			return err
		}
		permissionMap[perm.Code] = *perm
	}

	for _, perm := range permissionMap {
		if err := ensureRolePermission(db, superAdminRole.ID, perm.ID); err != nil {
			return err
		}
	}

	editorPermissionCodes := []string{
		"article:read",
		"article:create",
		"article:update",
		"category:read",
		"category:create",
		"category:update",
		"tag:read",
		"tag:create",
		"tag:update",
		"user:me",
		"user:logout",
	}
	for _, code := range editorPermissionCodes {
		perm, ok := permissionMap[code]
		if !ok {
			continue
		}
		if err := ensureRolePermission(db, editorRole.ID, perm.ID); err != nil {
			return err
		}
	}

	adminUser, err := ensureAdminUser(db)
	if err != nil {
		return err
	}
	return ensureUserRole(db, adminUser.ID, superAdminRole.ID)
}

func ensureAdminUser(db *gorm.DB) (*User, error) {
	user := &User{}
	if err := db.Where("username = ?", "admin").First(user).Error; err == nil {
		return user, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user = &User{
		Username:     "admin",
		PasswordHash: string(passwordHash),
		Nickname:     "Loe",
		Type:         "admin",
		Status:       "active",
	}
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func ensureRole(db *gorm.DB, name, code, description string) (*Role, error) {
	role := &Role{}
	if err := db.Where("code = ?", code).First(role).Error; err == nil {
		return role, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	role = &Role{Name: name, Code: code, Description: description}
	if err := db.Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func ensurePermission(db *gorm.DB, item Permission) (*Permission, error) {
	perm := &Permission{}
	if err := db.Where("code = ?", item.Code).First(perm).Error; err == nil {
		return perm, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	perm = &Permission{
		Name:        item.Name,
		Code:        item.Code,
		Resource:    item.Resource,
		Action:      item.Action,
		Description: item.Description,
	}
	if err := db.Create(perm).Error; err != nil {
		return nil, err
	}
	return perm, nil
}

func ensureRolePermission(db *gorm.DB, roleID, permissionID int64) error {
	var count int64
	if err := db.Model(&RolePermission{}).Where("role_id = ? AND permission_id = ?", roleID, permissionID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return db.Create(&RolePermission{RoleID: roleID, PermissionID: permissionID}).Error
}

func ensureUserRole(db *gorm.DB, userID, roleID int64) error {
	var count int64
	if err := db.Model(&UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return db.Create(&UserRole{UserID: userID, RoleID: roleID}).Error
}
