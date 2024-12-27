package authentication

import "modular_chassis/echo/pkg/errorz"

var (
	Unauthorized errorz.Type = 10001
	UserNotFound errorz.Type = 10002
)
