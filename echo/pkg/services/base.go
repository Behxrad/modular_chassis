package services

type BaseReq struct {
	Mobile string `json:"mobile"`
}

type BaseResp struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
