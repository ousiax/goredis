// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis

type resptransaction interface {
	// DISCARD
	// Discard all commands issued after MULTI
	// Simple string reply: always OK.
	Discard() (string, error)

	// EXEC
	// Execute all commands issued after MULTI
	// Array reply: each element being the reply to each of the commands in the atomic transaction.
	// When using WATCH, EXEC can return a Null reply if the execution was aborted.
	Exec() ([]interface{}, error)

	// MULTI
	// Mark the start of a transaction block
	// Simple string reply: always OK.
	Multi() (string, error)

	// UNWATCH
	// Forget about all watched keys
	// Simple string reply: always OK.
	UnWatch() (string, error)

	// WATCH key [key ...]
	// Watch the given keys to determine execution of the MULTI/EXEC block
	// Simple string reply: always OK.
	Watch(key interface{}, keys ...interface{}) (string, error)
}
