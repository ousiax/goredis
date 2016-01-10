package redis

type RESPStrings interface {
	// APPEND key value Append a value to a key
	//     Append(key, value string) (int, error)
	// BITCOUNT key [start end] Count set bits in a string
	// BITOP operation destkey key [key ...] Perform bitwise operations between strings
	// BITPOS key bit [start] [end] Find first bit set or clear in a string
	// DECR key Decrement the integer value of a key by one
	//     Decr(key string) (int, error)
	// DECRBY key decrement Decrement the integer value of a key by the given number
	// GET key Get the value of a key
	Get(key string) (interface{}, error)
	// GETBIT key offset Returns the bit value at offset in the string value stored at key
	// GETRANGE key start end Get a substring of the string stored at a key
	// GETSET key value Set the string value of a key and return its old value
	// INCR key Increment the integer value of a key by one
	//     Incr(key string) (int, error)
	// INCRBY key increment Increment the integer value of a key by the given amount
	// INCRBYFLOAT key increment Increment the float value of a key by the given amount
	// MGET key [key ...] Get the values of all the given keys
	// MSET key value [key value ...] Set multiple keys to multiple values
	// MSETNX key value [key value ...] Set multiple keys to multiple values, only if none of the keys exist
	// PSETEX key milliseconds value Set the value and expiration in milliseconds of a key
	// SET key value [EX seconds] [PX milliseconds] [NX|XX] Set the string value of a key
	//     Set(key string, args ...interface{}) (string, error)
	// SETBIT key offset value Sets or clears the bit at offset in the string value stored at key
	// SETEX key seconds value Set the value and expiration of a key
	// SETNX key value Set the value of a key, only if the key does not exist
	// SETRANGE key offset value Overwrite part of a string at key starting at the specified offset
	// STRLEN key Get the length of the value stored in a key
}
