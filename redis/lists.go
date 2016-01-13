package redis

type respLists interface {
	// BLPOP key [key ...] timeout
	// Remove and get the first element in a list, or block until one is available

	// BRPOP key [key ...] timeout
	// Remove and get the last element in a list, or block until one is available

	// BRPOPLPUSH source destination timeout
	// Pop a value from a list, push it to another list and return it; or block until one is available

	// LINDEX key index
	// Get an element from a list by its index

	// LINSERT key BEFORE|AFTER pivot value
	// Insert an element before or after another element in a list

	// LLEN key
	// Get the length of a list

	// LPOP key
	// Remove and get the first element in a list

	// LPUSH key value [value ...]
	// Prepend one or multiple values to a list

	// LPUSHX key value
	// Prepend a value to a list, only if the list exists

	// LRANGE key start stop
	// Get a range of elements from a list

	// LREM key count value
	// Remove elements from a list

	// LSET key index value
	// Set the value of an element in a list by its index

	// LTRIM key start stop
	// Trim a list to the specified range

	// RPOP key
	// Remove and get the last element in a list

	// RPOPLPUSH source destination
	// Remove the last element in a list, prepend it to another list and return it

	// RPUSH key value [value ...]
	// Append one or multiple values to a list

	// RPUSHX key value
	// Append a value to a list, only if the list exists
}
