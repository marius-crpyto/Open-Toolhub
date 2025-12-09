package server

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Ok(data any) Response {
	return Response{Code: 0, Data: data}
}

func Err(code int, msg string) Response {
	return Response{Code: code, Message: msg}
}
