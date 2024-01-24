package utility

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	infra_error "github.com/zaza-hikayat/go-fiber/internal/infrastructure/errors"
)

type IResponseClient interface {
	JSON(c *fiber.Ctx, data interface{}, meta *Meta) error
	HttpError(c *fiber.Ctx, err error) error
	Message(msg string) *responseClient
	SetStatus(code int) *responseClient
}

type responseClient struct {
	statusCode int
	message    string
}

type (
	ResponseMessage struct {
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
		Meta    *Meta       `json:"meta,omitempty"`
	}
	Meta struct {
		Page      int     `json:"page,omitempty"`
		PerPage   int     `json:"perPage,omitempty"`
		TotalPage float64 `json:"totalPage"`
		Total     int     `json:"total"`
	}
)

func NewResponseClient() IResponseClient {
	return &responseClient{}
}

func (*responseClient) HttpError(c *fiber.Ctx, err error) error {
	var respError infra_error.HttpError

	if cerr, ok := err.(*infra_error.CommonError); ok {
		respError = cerr.ToHttpError()
	} else {
		respError = infra_error.NewCommonError(infra_error.UNKNOWN_ERROR, err).ToHttpError()
	}

	c.Status(respError.HttpStatusNumber).JSON(respError)
	return nil
}
func (r *responseClient) Message(msg string) *responseClient {
	r.message = msg
	return r
}
func (r *responseClient) SetStatus(code int) *responseClient {
	r.statusCode = code
	return r
}

func (r *responseClient) JSON(c *fiber.Ctx, data interface{}, meta *Meta) error {
	statusCode := http.StatusOK
	if r.statusCode != 0 {
		statusCode = r.statusCode
	}
	response := ResponseMessage{
		Message: r.message,
		Data:    data,
		Meta:    meta,
	}
	return c.Status(statusCode).JSON(response)
}
