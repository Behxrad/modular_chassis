package service

import (
	"context"
	"math/rand"
	"modular_chassis/echo/pkg/services/finance"
)

func (Base) GetBalance(context.Context, finance.GetBalanceRequest) (finance.GetBalanceResponse, error) {
	return finance.GetBalanceResponse{Amount: rand.Int63n(10_000_000)}, nil
}
