package model

type UploadResponse struct {
	ID           int64  `json:"id"`
	OriginalName string `json:"original_name"`
	URL          string `json:"url"`
	Path         string `json:"path"`
	ContentType  string `json:"content_type"`
	Size         int64  `json:"size"`
	BizType      string `json:"biz_type"`
}