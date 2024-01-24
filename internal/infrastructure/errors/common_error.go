package errors

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ErrorCode = uint

type CommonError struct {
	ClientMessage   string      `json:"message"`
	SystemMessage   interface{} `json:"data"`
	ErrorCode       ErrorCode   `json:"code"`
	ErrorStatusCode uint        `json:"-"`
	ErrorMessage    *string     `json:"-"`
	ErrorTrace      *string     `json:"-"`
}

func (err CommonError) Error() string {
	var (
		errMsg string
		errTrc string
	)

	if err.ErrorMessage != nil {
		errMsg = *err.ErrorMessage
	}

	if err.ErrorTrace != nil {
		errTrc = *err.ErrorTrace
	}

	return fmt.Sprintf("Error: %s. Trace: %s", errMsg, errTrc)
}

func NewCommonError(errCode ErrorCode, err error) *CommonError {
	var errMsg *string
	var errTrace *string
	var clientMessage string = "Unknown error."
	var systemMessage interface{} = "Unknown error."
	var commonError = ErrorDicts[errCode]
	var validationMessage interface{}

	if err != nil {
		s := err.Error()
		errMsg = &s

		ss := fmt.Sprintf("%+v", err)
		errTrace = &ss

		if errCode == UNKNOWN_ERROR {
			systemMessage = ss
		}

		if _err, ok := err.(validation.Errors); ok {
			validationMessage = _err
		}
	}

	if commonError == nil {
		return &CommonError{
			ClientMessage: clientMessage,
			SystemMessage: systemMessage,
			ErrorCode:     errCode,
			ErrorTrace:    errTrace,
			ErrorMessage:  errMsg,
		}
	}

	if _err, ok := err.(*CommonError); ok {
		commonError = _err

	}
	if validationMessage != nil {
		commonError.SystemMessage = validationMessage
	}
	return &CommonError{
		ClientMessage:   commonError.ClientMessage,
		SystemMessage:   commonError.SystemMessage,
		ErrorCode:       errCode,
		ErrorTrace:      errTrace,
		ErrorMessage:    errMsg,
		ErrorStatusCode: commonError.ErrorStatusCode,
	}
}

func (err *CommonError) SetSystemMessage(message interface{}) {
	if _err, ok := message.(validation.Errors); ok {
		err.SystemMessage = _err
	}
	err.SystemMessage = message
}
