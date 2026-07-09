package db

type Article struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	Title       string `gorm:"type:varchar(255);not null"`
	Slug        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	Summary     string `gorm:"type:text"`
	ContentMd   string `gorm:"type:longtext"`
	ContentHTML string `gorm:"type:longtext"`
	CoverImage  string `gorm:"type:varchar(255)"`
	CategoryID  int64  `gorm:"default:0"`
	TagIds      string `gorm:"type:text"`
	Status      string `gorm:"type:varchar(32);index;default:'draft'"`
	CreatedAt   string `gorm:"type:varchar(32)"`
	UpdatedAt   string `gorm:"type:varchar(32)"`
	PublishedAt string `gorm:"type:varchar(32)"`
}