package errors

import "net/http"

const (
	INVALID_USER    ErrorCode = 4_00_000_1
	INVALID_PAYLOAD ErrorCode = 4_00_000_2
	INVALID_DATA    ErrorCode = 4_00_000_3
	INVALID_TOKEN   ErrorCode = 4_00_000_4
	UNKNOWN_ERROR   ErrorCode = 5_00_000_1
)

var ErrorDicts = map[ErrorCode]*CommonError{
	INVALID_USER: {
		ClientMessage:   "Invalid user",
		SystemMessage:   "Invalid user",
		ErrorStatusCode: http.StatusForbidden,
	},
	UNKNOWN_ERROR: {
		ClientMessage:   "Unknown error",
		SystemMessage:   "Unknown error",
		ErrorStatusCode: http.StatusInternalServerError,
	},
	INVALID_PAYLOAD: {
		ClientMessage:   "Invalid payload",
		SystemMessage:   "Invalid payload",
		ErrorStatusCode: http.StatusBadRequest,
	},
	INVALID_DATA: {
		ClientMessage:   "Invalid data",
		SystemMessage:   "Invalid data",
		ErrorStatusCode: http.StatusNotFound,
	},
	INVALID_TOKEN: {
		ClientMessage:   "Invalid token value",
		SystemMessage:   "Invalid token value",
		ErrorStatusCode: http.StatusUnauthorized,
	},
}
