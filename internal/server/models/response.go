package models

import "strings"

type Response struct {
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta"`
	Message string      `json:"message" example:"Ошибка обработки"`
}

func (r *Response) NewWithMessage(d interface{}, m string) Response {
	return Response{
		Data:    d,
		Message: strings.Trim(m, " "),
	}
}

func (r *Response) ErrorResponse(err string) Response {
	return Response{
		Message: strings.Trim(err, " "),
	}
}
