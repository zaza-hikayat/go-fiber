package errors

import (
	"net/http"
)

type HttpError struct {
	CommonError
	HttpStatusNumber int    `json:"-"`
	HttpStatusName   string `json:"type"`
}

func (err HttpError) Error() string {
	return err.ClientMessage
}

func (err CommonError) ToHttpError() HttpError {
	httpStatusNumber := http.StatusInternalServerError
	if err.ErrorStatusCode != 0 {
		httpStatusNumber = int(err.ErrorStatusCode)
	}

	return HttpError{
		CommonError:      err,
		HttpStatusNumber: httpStatusNumber,
		HttpStatusName:   GetHttpStatusText(httpStatusNumber),
	}
}
