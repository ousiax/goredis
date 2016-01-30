// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis

// MakeSlice make a new slice: [ p[0],p[1],...,p[n],opt[0],opt[1],...,opt[n] ]
func MakeSlice(opt []interface{}, p ...interface{}) []interface{} {
	l := len(opt) + len(p)
	a := make([]interface{}, 0, l)
	for _, v := range p {
		a = append(a, v)
	}
	for _, v := range opt {
		a = append(a, v)
	}
	return a
}
