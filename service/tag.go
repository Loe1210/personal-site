package service

import (
	"context"

	mysqlDriver "github.com/go-sql-driver/mysql"

	tagmodel "github.com/Loe1210/personal-site/biz/model/tag"
	dbmodel "github.com/Loe1210/personal-site/dal/db"
	"github.com/Loe1210/personal-site/pkg/errno"
)

func toTagModel(item *dbmodel.Tag) *tagmodel.Tag {
	if item == nil {
		return nil
	}

	return &tagmodel.Tag{
		ID:          item.ID,
		Name:        item.Name,
		Slug:        item.Slug,
		Description: item.Description,
		CreatedAt:   formatTime(item.CreatedAt),
		UpdatedAt:   formatTime(item.UpdatedAt),
	}
}

func CreateTag(_ context.Context, req *tagmodel.CreateTagRequest) (*tagmodel.CreateTagResponse, error) {
	record := &dbmodel.Tag{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := dbmodel.DB.Create(record).Error; err != nil {
		if mysqlErr, ok := err.(*mysqlDriver.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, errno.TagConflict
		}
		return nil, errno.Internal
	}

	return &tagmodel.CreateTagResponse{
		Tag:     toTagModel(record),
		Message: "tag created",
	}, nil
}

func UpdateTag(_ context.Context, req *tagmodel.UpdateTagRequest) (*tagmodel.UpdateTagResponse, error) {
	var record dbmodel.Tag
	if err := dbmodel.DB.First(&record, req.ID).Error; err != nil {
		return nil, errno.TagNotFound
	}

	record.Name = req.Name
	record.Slug = req.Slug
	record.Description = req.Description

	if err := dbmodel.DB.Save(&record).Error; err != nil {
		if mysqlErr, ok := err.(*mysqlDriver.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, errno.TagConflict
		}
		return nil, errno.Internal
	}

	return &tagmodel.UpdateTagResponse{
		Tag:     toTagModel(&record),
		Message: "tag updated",
	}, nil
}

func DeleteTag(_ context.Context, req *tagmodel.DeleteTagRequest) (*tagmodel.DeleteTagResponse, error) {
	var record dbmodel.Tag
	if err := dbmodel.DB.First(&record, req.ID).Error; err != nil {
		return nil, errno.TagNotFound
	}

	var articleTagCount int64
	if err := dbmodel.DB.Model(&dbmodel.ArticleTag{}).Where("tag_id = ?", req.ID).Count(&articleTagCount).Error; err != nil {
		return nil, errno.Internal
	}
	if articleTagCount > 0 {
		return nil, errno.TagInUse
	}

	if err := dbmodel.DB.Delete(&record).Error; err != nil {
		return nil, errno.Internal
	}

	return &tagmodel.DeleteTagResponse{Message: "tag deleted"}, nil
}

func ListTags(_ context.Context, _ *tagmodel.ListTagsRequest) (*tagmodel.ListTagsResponse, error) {
	var records []dbmodel.Tag

	if err := dbmodel.DB.Order("id DESC").Find(&records).Error; err != nil {
		return nil, err
	}

	list := make([]*tagmodel.Tag, 0, len(records))
	for i := range records {
		list = append(list, toTagModel(&records[i]))
	}

	return &tagmodel.ListTagsResponse{
		List: list,
	}, nil
}
