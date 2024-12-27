package authentication

import "modular_chassis/echo/pkg/errs"

var (
	Unauthorized errs.ResponseCode = 40002
	UserNotFound errs.ResponseCode = 40003
)
