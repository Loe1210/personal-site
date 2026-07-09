package category

import (
	"context"
	"time"

	dbmodel "github.com/Loe1210/personal-site/biz/dal/db"
	categorymodel "github.com/Loe1210/personal-site/biz/model/category"
)

const timeLayout = "2006-01-02 15:04:05"

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format(timeLayout)
}

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
		return nil, err
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