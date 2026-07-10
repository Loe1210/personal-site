package service

import (
	"context"
	
	mysqlDriver "github.com/go-sql-driver/mysql"

	dbmodel "github.com/Loe1210/personal-site/dal/db"
	categorymodel "github.com/Loe1210/personal-site/biz/model/category"
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

func ListCategories(_ context.Context, _ *categorymodel.ListCategoriesRequest) (*categorymodel.ListCategoriesResponse, error) {
	var records []dbmodel.Category

	if err := dbmodel.DB.Order("id DESC").Find(&records).Error; err != nil {
		return nil, err
	}

	list := make([]*categorymodel.Category, 0, len(records))
	for i := range records {
		list = append(list, toCategoryModel(&records[i]))
	}

	return &categorymodel.ListCategoriesResponse{
		List: list,
	}, nil
}



