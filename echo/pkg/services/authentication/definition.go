package authentication

import (
	"context"
	"modular_chassis/echo/pkg/services"
)

type UserResponse struct {
	Username string `json:"username" binding:"required,min=5,max=20" swaggo:"required,min=5,max=20"`
	Email    string `json:"email" binding:"required,email" swaggo:"required,email"`
	Age      int    `json:"age" binding:"gte=18" swaggo:"minimum=18"`
}

type UserRequest struct {
	ID int `json:"id" binding:"required"`
}

type TestRequest struct {
	services.BaseReq
	Input string `json:"input" binding:"required"`
}

type TestResponse struct {
	services.BaseResp
	Output string `json:"output" binding:"required"`
}

type (
	Service interface {
		Test(ctx context.Context, request TestRequest) (TestResponse, error)
		FetchUserID(ctx context.Context, request UserRequest) (UserResponse, error)
	}
	Service2 interface {
		Test(ctx context.Context, request TestRequest) (TestResponse, error)
	}
)
