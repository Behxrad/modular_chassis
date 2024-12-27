package services

import "modular_chassis/echo/pkg/errs"

var (
	GeneralError  errs.ResponseCode = 10000
	InternalError errs.ResponseCode = 20000
	UnknownError  errs.ResponseCode = 30000
	BadRequest    errs.ResponseCode = 40000
)
