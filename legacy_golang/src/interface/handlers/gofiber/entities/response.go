package gofiberentities

import (
	"github.com/gofiber/fiber/v2"
)

type IResponse interface {
	Success(code int, data any) IResponse
	Error(code int, tractId, msg string) IResponse
	Response() error
}

type Response struct {
	StatusCode    int
	Data          any
	ErrorResponse *ErrorResponseType
	Context       *fiber.Ctx
	IsError       bool
}

type ErrorResponseType struct {
	TraceId string `json:"trace_id"`
	Msg     string `json:"message"`
}

func NewResponse(c *fiber.Ctx) IResponse {
	return &Response{
		Context: c,
	}
}

func (r *Response) Success(code int, data any) IResponse {
	r.StatusCode = code
	r.Data = data
	// kawaiilogger.InitKawaiiLogger(r.Context, &r.Data, code).Print()
	return r
}
func (r *Response) Error(code int, tractId, msg string) IResponse {
	r.StatusCode = code
	r.ErrorResponse = &ErrorResponseType{
		TraceId: tractId,
		Msg:     msg,
	}
	r.IsError = true
	// kawaiilogger.InitKawaiiLogger(r.Context, &r.ErrorRes, code).Print()
	return r
}
func (r *Response) Response() error {
	return r.Context.Status(r.StatusCode).JSON(func() any {
		if r.IsError {
			return &r.ErrorResponse
		}
		return &r.Data
	}())
}

type PaginateRes struct {
	Data      any `json:"data"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalPage int `json:"total_page"`
	TotalItem int `json:"total_item"`
}
