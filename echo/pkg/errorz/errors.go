package errorz

type Type int32

type ServiceError struct {
	ErrorType    Type
	Cause        error
	TraceMessage string
}

var (
	GeneralError  Type = 10000
	InternalError Type = 20000
	UnknownError  Type = 30000
)
