package response

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) Body {
	return Body{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func Error(code int, message string) Body {
	return Body{
		Code:    code,
		Message: message,
	}
}