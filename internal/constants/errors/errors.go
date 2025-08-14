package errors

import (
	"errors"
	"net/http"
)

var (
	ErrUnexpected                  = errors.New("unexpected error")
	ErrInternalServerError         = errors.New("internal server error")
	ErrBeneficiaryAlreadyExist     = errors.New("beneficiary already exist")
	ErrBeneficiaryNotFound         = errors.New("beneficiary not found")
	ErrInividualsNotFound          = errors.New("one or more beneficiaries not found")
	ErrGroupBeneficaryAlreadyExist = errors.New("group beneficiary already exist")
	ErrRequestTimeout              = errors.New("request timeout")
	ErrAccountNotFound             = errors.New("account not found")
	ErrBadRequest                  = errors.New("bad request")
	ErrInvalidData                 = errors.New("invalid data")
	ErrUnauthorized                = errors.New("unauthorized")
	ErrActionNotAllowed            = errors.New("action not allowed")
	ErrInvalidPhoneNumber          = errors.New("invalid phone number")
)

var ErrorMap = map[error]int{
	ErrUnexpected:                  http.StatusInternalServerError,
	ErrInternalServerError:         http.StatusInternalServerError,
	ErrBeneficiaryAlreadyExist:     http.StatusBadRequest,
	ErrBeneficiaryNotFound:         http.StatusNotFound,
	ErrInividualsNotFound:          http.StatusNotFound,
	ErrGroupBeneficaryAlreadyExist: http.StatusBadRequest,
	ErrRequestTimeout:              http.StatusRequestTimeout,
	ErrAccountNotFound:             http.StatusNotFound,
	ErrBadRequest:                  http.StatusBadRequest,
	ErrInvalidData:                 http.StatusBadRequest,
	ErrUnauthorized:                http.StatusUnauthorized,
	ErrActionNotAllowed:            http.StatusForbidden,
	ErrInvalidPhoneNumber:          http.StatusBadRequest,
}
