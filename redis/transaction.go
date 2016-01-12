// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package redis

type respTransaction interface {
	// DISCARD Discard all commands issued after MULTI
	// EXEC Execute all commands issued after MULTI
	// MULTI Mark the start of a transaction block
	// UNWATCH Forget about all watched keys
	// WATCH key [key ...] Watch the given keys to determine execution of the MULTI/EXEC block
}
