// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis

type resplists interface {
	// BLPOP key [key ...] timeout
	// Remove and get the first element in a list, or block until one is available
	// Array reply: specifically:
	//      A nil multi-bulk when no element could be popped and the timeout expired.
	//      A two-element multi-bulk with the first element being the name of the key where an element was popped and the second element being the value of the popped element.
	BLPop(key interface{}, timeout int, keys ...interface{}) ([]interface{}, error)

	// BRPOP key [key ...] timeout
	// Remove and get the last element in a list, or block until one is available
	// Array reply: specifically:
	//      A nil multi-bulk when no element could be popped and the timeout expired.
	//      A two-element multi-bulk with the first element being the name of the key where an element was popped and the second element being the value of the popped element.
	BRPop(key interface{}, timeout int, keys ...interface{}) (int, error)

	// BRPOPLPUSH source destination timeout
	// Pop a value from a list, push it to another list and return it; or block until one is available
	// Bulk string reply: the element being popped from source and pushed to destination. If timeout is reached, a Null reply is returned.
	BRPopLPush(source, destination interface{}, timeout int) (interface{}, error)

	// LINDEX key index
	// Get an element from a list by its index
	// Bulk string reply: the requested element, or nil when index is out of range.
	LIndex(key interface{}, index int) (interface{}, error)

	// LINSERT key BEFORE|AFTER pivot value
	// Insert an element before or after another element in a list
	// Integer reply: the length of the list after the insert operation, or -1 when the value pivot was not found.
	LInsert(key interface{}, option string, pivot, value interface{}) (int, error)

	// LLEN key
	// Get the length of a list
	// Integer reply: the length of the list at key.
	LLen(key interface{}) (int, error)

	// LPOP key
	// Remove and get the first element in a list
	// Bulk string reply: the value of the first element, or nil when key does not exist.
	LPop(key interface{}) (interface{}, error)

	// LPUSH key value [value ...]
	// Prepend one or multiple values to a list
	// Integer reply: the length of the list after the push operations.
	LPush(key, value interface{}, values ...interface{}) (int, error)

	// LPUSHX key value
	// Prepend a value to a list, only if the list exists
	// Integer reply: the length of the list after the push operation.
	LPushx(key, value interface{}) (int, error)

	// LRANGE key start stop
	// Get a range of elements from a list
	// Array reply: list of elements in the specified range.
	LRange(key interface{}, start, stop int) (int, error)

	// LREM key count value
	// Remove elements from a list
	// The count argument influences the operation in the following ways:
	//      count > 0: Remove elements equal to value moving from head to tail.
	//      count < 0: Remove elements equal to value moving from tail to head.
	//      count = 0: Remove all elements equal to value.
	// Integer reply: the number of removed elements.
	LRem(key interface{}, count int, value interface{}) (int, error)

	// LSET key index value
	// Set the value of an element in a list by its index
	// Simple string reply
	LSet(key interface{}, index int, value interface{}) (string, error)

	// LTRIM key start stop
	// Trim a list to the specified range
	// Simple string reply
	LTrim(key interface{}, start, stop int) (string, error)

	// RPOP key
	// Remove and get the last element in a list
	// Bulk string reply: the value of the last element, or nil when key does not exist.
	RPop(key interface{}) (interface{}, error)

	// RPOPLPUSH source destination
	// Remove the last element in a list, prepend it to another list and return it
	// Bulk string reply: the element being popped and pushed.
	RPopLPush(source, destination interface{}) (interface{}, error)

	// RPUSH key value [value ...]
	// Append one or multiple values to a list
	// Integer reply: the length of the list after the push operation.
	RPush(key, value interface{}, values ...interface{})

	// RPUSHX key value
	// Append a value to a list, only if the list exists
	// Integer reply: the length of the list after the push operation.
	RPushx(key, value interface{}, values ...interface{}) (int, error)
}
