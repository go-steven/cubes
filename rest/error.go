package rest

import (
	"fmt"
)

type ErrorCode uint

const (
	BADREQUEST_ERROR     ErrorCode = 400
	INTERNAL_ERROR       ErrorCode = 500
	NOTFOUND_ERROR       ErrorCode = 404
	UNAUTHORIZED_ERROR   ErrorCode = 401
	INVALID_PASSWD_ERROR ErrorCode = 409
)

type APIError struct {
	Code ErrorCode `json:"code,omitempty"`
	Msg  string    `json:"msg,omitempty"`
}

func (this APIError) Error() string {
	return fmt.Sprintf("CODE:%d, MSG:%s", this.Code, this.Msg)
}
