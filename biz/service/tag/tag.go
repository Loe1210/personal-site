package tag

import (
	"context"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"

	dbmodel "github.com/Loe1210/personal-site/biz/dal/db"
	tagmodel "github.com/Loe1210/personal-site/biz/model/tag"
	"github.com/Loe1210/personal-site/pkg/errno"
)

const timeLayout = "2006-01-02 15:04:05"

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Local().Format(timeLayout)
}

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
