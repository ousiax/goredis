// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis

import (
	"errors"
	"fmt"
	"strconv"
)

// Int parses a RESP Integer to int.
func Int(p interface{}) (int, error) {
	b, e := p.([]byte)
	if e {
		return strconv.Atoi(string(b))
	}
	return 0, strconv.ErrRange
}

// Float64 parses a RESP Bulk String to a float64 number.
func Float64(p interface{}) (float64, error) {
	b, e := p.([]byte)
	if e {
		return strconv.ParseFloat(string(b), 64)
	}
	return 0, strconv.ErrRange
}

// Stringx parses a RESP Bulk String or a Simple String to a string or a nil, otherwise a nil when a error occured.
func Stringx(p interface{}) (interface{}, error) {
	switch v := p.(type) {
	case []byte:
		return string(v), nil
	case nil:
		return nil, nil
	default:
		return nil, errors.New("redis.Stringx(interface{}): Protocol error.")
	}
}

// String parses RESP Bulk String or Simple String to a string, otherwise a empty string.
// usually, the p is a string or a nil (i.e. a zero value).
func String(p interface{}) (string, error) {
	v, e := Stringx(p)
	s, _ := v.(string)
	return s, e
}

// Strings parses a RESP reply (array apply) to a string arrray (may contain null value).
func Strings(p interface{}) ([]interface{}, error) {
	rsp, e := p.([]interface{})
	if !e {
		return nil, errors.New(fmt.Sprintf("redis.Strings(interface{}): interface conversion, interface is %T, not []interface{}.", p))
	}
	for i, v := range rsp {
		p, _ := Stringx(v) // Error check should be not required for performance.
		rsp[i] = p
	}
	return rsp, nil
}
