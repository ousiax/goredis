// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis

import (
	"errors"
	"fmt"
	"strconv"
)

func parseInt(p interface{}) (int, error) {
	v, e := p.(int)
	if e {
		return v, nil
	}
	return v, strconv.ErrRange
}

// parseFloat parses a RESP Bulk String to a float64 number.
func parseFloat(p interface{}) (float64, error) {
	switch v := p.(type) {
	case []byte:
		return strconv.ParseFloat(string(v), 64)
	default:
		return 0.0, strconv.ErrRange
	}
}

// parseStringEx parses a RESP Bulk String or a Simple String to a string or a nil, otherwise a nil when a error occured.
func parseStringEx(p interface{}) (interface{}, error) {
	switch v := p.(type) {
	case []byte:
		return string(v), nil
	case nil:
		return nil, nil
	case string:
		return v, nil
	default:
		return nil, errors.New(fmt.Sprintf("redis.parseStringEx: interface conversion, interface is %T, not string.", p))
	}
}

// parseString returns a string if p is a string type, otherwise a empty string.
// usually, the p is a string or a nil (i.e. a zero value).
func parseString(p interface{}) (string, error) {
	rsp, err := parseStringEx(p)
	s, _ := rsp.(string)
	return s, err
}

// parseStrings parses a RESP reply (array apply) to a string arrray (may contain null value).
func parseStrings(p interface{}) ([]interface{}, error) {
	rsp, e := p.([]interface{})
	if !e {
		return nil, errors.New(fmt.Sprintf("redis.parseStrings: interface conversion, interface is %T, not []interface{}.", p))
	}
	for i, v := range rsp {
		p, _ := parseStringEx(v) // Error check should be not required for performance.
		rsp[i] = p
	}
	return rsp, nil
}
