package services

type ServiceRequest[P any] struct {
	Payload P
}

type ServiceResponse[P any] struct {
	Payload P
}
