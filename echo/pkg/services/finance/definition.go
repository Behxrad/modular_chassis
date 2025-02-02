package finance

import "context"

type GetBalanceRequest struct {
	UserID string `json:"userID"`
}

type GetBalanceResponse struct {
	Amount int64 `json:"amount"`
}

type TestRequest struct {
	Value string `json:"value"`
}

type TestResponse struct {
	Value string `json:"value"`
}

type Service interface {
	GetBalance(context.Context, GetBalanceRequest) (GetBalanceResponse, error)
	Test(context.Context, TestRequest) (TestResponse, error)
}
