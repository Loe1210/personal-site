package tag

type UpdateTagRequest struct {
	ID          int64  `path:"id" json:"id" form:"id" query:"id"`
	Name        string `json:"name" form:"name"`
	Slug        string `json:"slug" form:"slug"`
	Description string `json:"description" form:"description"`
}

type UpdateTagResponse struct {
	Tag     *Tag   `json:"tag"`
	Message string `json:"message"`
}

type DeleteTagRequest struct {
	ID int64 `path:"id" json:"id" form:"id" query:"id"`
}

type DeleteTagResponse struct {
	Message string `json:"message"`
}
