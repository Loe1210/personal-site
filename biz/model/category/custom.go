package category

type GetCategoryRequest struct {
	ID uint `path:"id" json:"id" form:"id" query:"id"`
}

type GetCategoryResponse struct {
	Category *Category `json:"category"`
}

type UpdateCategoryRequest struct {
	ID          int64  `path:"id" json:"id" form:"id" query:"id"`
	Name        string `json:"name" form:"name"`
	Slug        string `json:"slug" form:"slug"`
	Description string `json:"description" form:"description"`
}

type UpdateCategoryResponse struct {
	Category *Category `json:"category"`
	Message  string    `json:"message"`
}

type DeleteCategoryRequest struct {
	ID int64 `path:"id" json:"id" form:"id" query:"id"`
}

type DeleteCategoryResponse struct {
	Message string `json:"message"`
}
