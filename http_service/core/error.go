package core

const (
	ErrCodeInternalError = 1000
	ErrCodeInvalidParams = 1001
	ErrCodeNotFound      = 1002
	ErrCodeAlreadyExists = 1003
	ErrCodeUnauthorized  = 1004
	ErrCodeForbidden     = 1005

	FailedToDeleteNode = 1006
	FailedToGetNode    = 1007
	FailedToAddNode    = 1008
)

type KnownError struct {
	Code int
	Data interface{}
	Msg  string
}

func NewKnownError(code int, data interface{}, msg string) *KnownError {
	e := &KnownError{}
	e.Code = code
	e.Data = data
	e.Msg = msg
	return e
}

func (e *KnownError) Error() string {
	return e.Msg
}
