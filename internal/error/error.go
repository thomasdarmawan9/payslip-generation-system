package error

import "net/http"

const (
	DefaultErrorCaseCode = "00"
)

const (
	BadRequest          = "Bad Request"
	InternalServerError = "internal server error"
	InvalidFormat       = "Invalid Format"
	InvalidMandatory    = "Invalid Mandatory"
	ErrUnauthorized     = "Unauthorized"
	NotFoundError       = "Not Found"
	ConflictError       = "Conflict"
)

var (
	ErrorMapCaseCode = map[string]string{
		BadRequest:          "00",
		InternalServerError: "00",
		InvalidFormat:       "01",
		InvalidMandatory:    "02",
		ErrUnauthorized:     "01",
		NotFoundError:       "03",
		ConflictError:       "04",
	}

	ErrorMapMessage = map[string]string{
		BadRequest:          "Bad Request",
		InternalServerError: "internal server error",
		InvalidFormat:       "Format %v tidak sesuai",
		InvalidMandatory:    "Field %v tidak boleh kosong",
		ErrUnauthorized:     "Unauthorized",
		NotFoundError:       "Data tidak ditemukan",
		ConflictError:       "Data sudah ada, Reason: %v",
	}

	ErrorMapHttpCode = map[string]int{
		BadRequest:          http.StatusBadRequest,
		InternalServerError: http.StatusInternalServerError,
		InvalidFormat:       http.StatusBadRequest,
		InvalidMandatory:    http.StatusBadRequest,
		ErrUnauthorized:     http.StatusUnauthorized,
		NotFoundError:       http.StatusNotFound,
		ConflictError:       http.StatusConflict,
	}
)
