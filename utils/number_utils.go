package utils

import (
	"fmt"
	"strconv"
)

func GreaterThan2interface(a, b interface{}) (bool, error) {
	aFloat64, err := Number2float64(a)
	if err != nil {
		return false, err
	}
	bFloat64, err := Number2float64(b)
	if err != nil {
		return false, err
	}
	return aFloat64 > bFloat64, nil
}

func Number2float64(a interface{}) (float64, error) {
	switch a := a.(type) {
	case int:
		return float64(a), nil
	case int32:
		return float64(a), nil
	case int64:
		return float64(a), nil
	case float64:
		return a, nil
	case float32:
		return float64(a), nil
	case uint:
		return float64(a), nil
	case uint32:
		return float64(a), nil
	case uint64:
		return float64(a), nil
	case string:
		return strconv.ParseFloat(a, 64)
	default:
		return 0, fmt.Errorf("unsupported type %T", a)
	}
}

func SubInterface(a, b interface{}) (interface{}, error) {
	aFloat64, err := Number2float64(a)
	if err != nil {
		return 0, err
	}
	bFloat64, err := Number2float64(b)
	if err != nil {
		return 0, err
	}
	return aFloat64 - bFloat64, nil
}

