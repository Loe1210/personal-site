package errno

type AppError struct {
	Code       int
	Message    string
	HTTPStatus int
}

func (e *AppError) Error() string {
	return e.Message
}

const (
	SuccessCode = 0
	ErrorCode   = 10000
)

var (
	Success = &AppError{Code: SuccessCode, Message: "success", HTTPStatus: 200}

	BadRequest   = &AppError{Code: 10001, Message: "request parameter invalid", HTTPStatus: 400}
	Unauthorized = &AppError{Code: 10002, Message: "unauthorized", HTTPStatus: 401}
	Forbidden    = &AppError{Code: 10003, Message: "forbidden", HTTPStatus: 403}
	NotFound     = &AppError{Code: 10004, Message: "resource not found", HTTPStatus: 404}
	Conflict     = &AppError{Code: 10005, Message: "resource conflict", HTTPStatus: 409}
	Internal     = &AppError{Code: 10006, Message: "internal server error", HTTPStatus: 500}

	CategoryNotFound = &AppError{Code: 20001, Message: "category not found", HTTPStatus: 400}
	ArticleNotFound  = &AppError{Code: 20002, Message: "article not found", HTTPStatus: 404}
	SlugConflict     = &AppError{Code: 20003, Message: "article slug already exists", HTTPStatus: 409}
	TagNotFound      = &AppError{Code: 20004, Message: "tag not found", HTTPStatus: 400}
	CategoryConflict = &AppError{Code: 20005, Message: "category already exists", HTTPStatus: 409}
	TagConflict      = &AppError{Code: 20006, Message: "tag already exists", HTTPStatus: 409}
	CategoryInUse    = &AppError{Code: 20007, Message: "category is used by articles", HTTPStatus: 409}
	TagInUse         = &AppError{Code: 20008, Message: "tag is used by articles", HTTPStatus: 409}
)
