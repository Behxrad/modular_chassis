package telecom

import "context"

type GetPackagesRequest struct {
	Operator string `json:"operator"`
}

type GetPackagesResponse struct {
	Packages []Package `json:"packages"`
}

type Package struct {
	Name string `json:"name"`
}

type TestRequest struct {
	Value string `json:"value"`
}

type TestResponse struct {
	Value string `json:"value"`
}

type Service interface {
	GetPackages(context.Context, GetPackagesRequest) (GetPackagesResponse, error)
	Test(context.Context, TestRequest) (TestResponse, error)
}
