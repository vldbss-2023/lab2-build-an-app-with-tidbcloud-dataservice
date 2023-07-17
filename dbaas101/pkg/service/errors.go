package service

import (
	"fmt"
	"net/http"
)

type ErrorInfo struct {
	Code    int    `json:"code"` // same as http code
	Message string `json:"message"`
}

func NewErrorInfo(statusCode int, message string) ErrorInfo {
	return ErrorInfo{Code: statusCode, Message: message}
}

func (err ErrorInfo) Error() string {
	return fmt.Sprintf("code: %d, message: %s", err.Code, err.Message)
}

func (err ErrorInfo) StatusCode() int {
	return err.Code
}

var (
	ErrInvalidParameter = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusBadRequest, msg) }
	ErrInternal         = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusInternalServerError, msg) }
	ErrUnauthorized     = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusBadRequest, msg) }
	ErrForbidden        = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusBadRequest, msg) }
	ErrNotSatisfied     = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusPreconditionRequired, msg) }
	ErrNotFound         = func(msg string) ErrorInfo { return NewErrorInfo(http.StatusBadRequest, msg) }
	ErrTooManyRequest   = func() ErrorInfo { return NewErrorInfo(http.StatusTooManyRequests, "Too many request.") }
	ErrNotImplemented   = func() ErrorInfo { return NewErrorInfo(http.StatusNotImplemented, "Not implemented.") }
)
