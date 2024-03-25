package utils

func Adr[T interface{}](s T) *T { return &s }
