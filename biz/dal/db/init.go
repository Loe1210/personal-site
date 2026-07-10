package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Loe1210/personal-site/pkg/configs"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

func Init() error {
	cfg := configs.LoadMySQLConfig()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := database.AutoMigrate(
		&User{},
		&Role{},
		&Permission{},
		&UserRole{},
		&RolePermission{},
		&Article{},
		&Category{},
		&Tag{},
		&ArticleTag{},
		&UploadFile{},
	); err != nil {
		return err
	}

	DB = database
	if err := seedDefaultUser(); err != nil {
		return err
	}
	if err := seedTestUser(); err != nil {
		return err
	}
	if err := seedRBAC(); err != nil {
		return err
	}

	return nil
}

func seedDefaultUser() error {
	var count int64
	if err := DB.Model(&User{}).Where("username = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Username:     "admin",
		PasswordHash: string(passwordHash),
		Nickname:     "Loe",
		Type:         "admin",
		Status:       "active",
	}

	return DB.Create(user).Error
}

func seedTestUser() error {
	var count int64
	if err := DB.Model(&User{}).Where("username = ?", "editor1").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Username:     "editor1",
		PasswordHash: string(passwordHash),
		Nickname:     "Editor",
		Type:         "admin",
		Status:       "active",
	}

	return DB.Create(user).Error
}

func seedRBAC() error {
	superAdminRole, err := ensureRole("Super Admin", "super_admin", "super administrator")
	if err != nil {
		return err
	}

	editorRole, err := ensureRole("Editor", "editor", "content editor")
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
		{Name: "Tag Read", Code: "tag:read", Resource: "tag", Action: "read", Description: "read tag"},
		{Name: "Tag Create", Code: "tag:create", Resource: "tag", Action: "create", Description: "create tag"},
		{Name: "User Me", Code: "user:me", Resource: "user", Action: "me", Description: "get current user"},
		{Name: "User Logout", Code: "user:logout", Resource: "user", Action: "logout", Description: "logout user"},
	}

	permissionMap := make(map[string]Permission)
	for _, item := range permissions {
		perm, err := ensurePermission(item)
		if err != nil {
			return err
		}
		permissionMap[perm.Code] = *perm
	}

	// super_admin 拥有全部权限
	for _, perm := range permissionMap {
		if err := ensureRolePermission(superAdminRole.ID, perm.ID); err != nil {
			return err
		}
	}

	// editor 只拥有部分权限
	editorPermissionCodes := []string{
		"article:read",
		"article:create",
		"article:update",
		"category:read",
		"tag:read",
		"user:me",
		"user:logout",
	}

	for _, code := range editorPermissionCodes {
		perm, ok := permissionMap[code]
		if !ok {
			continue
		}
		if err := ensureRolePermission(editorRole.ID, perm.ID); err != nil {
			return err
		}
	}

	var adminUser User
	if err := DB.Where("username = ?", "admin").First(&adminUser).Error; err != nil {
		return err
	}
	if err := ensureUserRole(adminUser.ID, superAdminRole.ID); err != nil {
		return err
	}

	var editorUser User
	if err := DB.Where("username = ?", "editor1").First(&editorUser).Error; err != nil {
		return err
	}
	if err := ensureUserRole(editorUser.ID, editorRole.ID); err != nil {
		return err
	}

	return nil
}

func ensureRole(name, code, description string) (*Role, error) {
	role := &Role{}
	if err := DB.Where("code = ?", code).First(role).Error; err == nil {
		return role, nil
	}

	role = &Role{
		Name:        name,
		Code:        code,
		Description: description,
	}
	if err := DB.Create(role).Error; err != nil {
		return nil, err
	}

	return role, nil
}

func ensurePermission(item Permission) (*Permission, error) {
	perm := &Permission{}
	if err := DB.Where("code = ?", item.Code).First(perm).Error; err == nil {
		return perm, nil
	}

	perm = &Permission{
		Name:        item.Name,
		Code:        item.Code,
		Resource:    item.Resource,
		Action:      item.Action,
		Description: item.Description,
	}
	if err := DB.Create(perm).Error; err != nil {
		return nil, err
	}

	return perm, nil
}

func ensureRolePermission(roleID, permissionID int64) error {
	var count int64
	if err := DB.Model(&RolePermission{}).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return DB.Create(&RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}).Error
}

func ensureUserRole(userID, roleID int64) error {
	var count int64
	if err := DB.Model(&UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	return DB.Create(&UserRole{
		UserID: userID,
		RoleID: roleID,
	}).Error
}