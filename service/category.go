package service

import (
	"context"

	mysqlDriver "github.com/go-sql-driver/mysql"

	categorymodel "github.com/Loe1210/personal-site/biz/model/category"
	dbmodel "github.com/Loe1210/personal-site/dal/db"
	"github.com/Loe1210/personal-site/pkg/errno"
)

func toCategoryModel(item *dbmodel.Category) *categorymodel.Category {
	if item == nil {
		return nil
	}

	return &categorymodel.Category{
		ID:          item.ID,
		Name:        item.Name,
		Slug:        item.Slug,
		Description: item.Description,
		CreatedAt:   formatTime(item.CreatedAt),
		UpdatedAt:   formatTime(item.UpdatedAt),
	}
}

func GetCategory(_ context.Context, req *categorymodel.GetCategoryRequest) (*categorymodel.GetCategoryResponse, error) {
	var record dbmodel.Category
	if err := dbmodel.DB.First(&record, req.ID).Error; err != nil {
		return nil, errno.CategoryNotFound
	}
	return &categorymodel.GetCategoryResponse{Category: toCategoryModel(&record)}, nil
}

func CreateCategory(_ context.Context, req *categorymodel.CreateCategoryRequest) (*categorymodel.CreateCategoryResponse, error) {
	record := &dbmodel.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := dbmodel.DB.Create(record).Error; err != nil {
		if mysqlErr, ok := err.(*mysqlDriver.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, errno.CategoryConflict
		}
		return nil, errno.Internal
	}

	return &categorymodel.CreateCategoryResponse{
		Category: toCategoryModel(record),
		Message:  "category created",
	}, nil
}

func UpdateCategory(_ context.Context, req *categorymodel.UpdateCategoryRequest) (*categorymodel.UpdateCategoryResponse, error) {
	var record dbmodel.Category
	if err := dbmodel.DB.First(&record, req.ID).Error; err != nil {
		return nil, errno.CategoryNotFound
	}

	record.Name = req.Name
	record.Slug = req.Slug
	record.Description = req.Description

	if err := dbmodel.DB.Save(&record).Error; err != nil {
		if mysqlErr, ok := err.(*mysqlDriver.MySQLError); ok && mysqlErr.Number == 1062 {
			return nil, errno.CategoryConflict
		}
		return nil, errno.Internal
	}

	return &categorymodel.UpdateCategoryResponse{
		Category: toCategoryModel(&record),
		Message:  "category updated",
	}, nil
}

func DeleteCategory(_ context.Context, req *categorymodel.DeleteCategoryRequest) (*categorymodel.DeleteCategoryResponse, error) {
	var record dbmodel.Category
	if err := dbmodel.DB.First(&record, req.ID).Error; err != nil {
		return nil, errno.CategoryNotFound
	}

	var articleCount int64
	if err := dbmodel.DB.Model(&dbmodel.Article{}).Where("category_id = ?", req.ID).Count(&articleCount).Error; err != nil {
		return nil, errno.Internal
	}
	if articleCount > 0 {
		return nil, errno.CategoryInUse
	}

	if err := dbmodel.DB.Delete(&record).Error; err != nil {
		return nil, errno.Internal
	}

	return &categorymodel.DeleteCategoryResponse{Message: "category deleted"}, nil
}

func ListCategories(_ context.Context, _ *categorymodel.ListCategoriesRequest) (*categorymodel.ListCategoriesResponse, error) {
	var records []dbmodel.Category

	if err := dbmodel.DB.Order("id DESC").Find(&records).Error; err != nil {
		return nil, errno.Internal
	}

	list := make([]*categorymodel.Category, 0, len(records))
	for i := range records {
		list = append(list, toCategoryModel(&records[i]))
	}

	return &categorymodel.ListCategoriesResponse{
		List: list,
	}, nil
}
