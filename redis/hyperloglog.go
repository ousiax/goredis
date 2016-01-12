// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package redis

type respHyperLogLog interface {
	// PFADD key element [element ...] Adds the specified elements to the specified HyperLogLog.
	// PFCOUNT key [key ...] Return the approximated cardinality of the set(s) observed by the HyperLogLog at key(s).
	// PFMERGE destkey sourcekey [sourcekey ...] Merge N different HyperLogLogs into a single one.
}
