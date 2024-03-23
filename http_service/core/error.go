package core


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


