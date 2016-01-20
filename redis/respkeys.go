// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis

type respkeys interface {
	// DEL key [key ...]
	// Delete a key
	// Integer reply: The number of keys that were removed.
	Del(key interface{}, keys ...interface{}) (int, error)

	// DUMP key
	// Return a serialized version of the value stored at the specified key.
	// Bulk string reply: the serialized value.
	Dump(key interface{}) (interface{}, error)

	// EXISTS key [key ...]
	// Determine if a key exists
	// Integer reply: The number of keys existing among the ones specified as arguments.
	Exists(key interface{}, keys ...interface{}) (int, error)

	// EXPIRE key seconds
	// Set a key's time to live in seconds
	// Integer reply, specifically:
	//     1 if the timeout was set.
	//     0 if key does not exist or the timeout could not be set.
	Expire(key interface{}, seconds int) (int, error)

	// EXPIREAT key timestamp
	// Set the expiration for a key as a UNIX timestamp
	// Integer reply, specifically:
	//     1 if the timeout was set.
	//     0 if key does not exist or the timeout could not be set.
	ExpireAt(key interface{}, timestamp int) (int, error)

	// KEYS pattern
	// Find all keys matching the given pattern
	// Array reply: list of keys matching pattern.
	Keys(pattern interface{}) ([]interface{}, error)

	// MIGRATE host port key destination-db timeout [COPY] [REPLACE]
	// Atomically transfer a key from a Redis instance to another one.
	// Simple string reply: The command returns OK on success.

	// MOVE key db
	// Move a key to another database
	// Integer reply, specifically:
	//     1 if key was moved.
	//     0 if key was not moved.
	Move(key interface{}, db int) (int, error)

	// OBJECT subcommand [arguments [arguments ...]]
	// Inspect the internals of Redis objects
	// The OBJECT command supports multiple sub commands:
	//     OBJECT REFCOUNT <key> returns the number of references of the value associated with the specified key. This command is mainly useful for debugging.
	//     OBJECT ENCODING <key> returns the kind of internal representation used in order to store the value associated with a key.
	//     OBJECT IDLETIME <key> returns the number of seconds since the object stored at the specified key is idle (not requested by read or write operations).
	//     While the value is returned in seconds the actual resolution of this timer is 10 seconds, but may vary in future implementations.
	// Objects can be encoded in different ways:
	//     Strings can be encoded as raw (normal string encoding) or int (strings representing integers in a 64 bit signed interval are encoded in this way in order to save space).
	//     Lists can be encoded as ziplist or linkedlist. The ziplist is the special representation that is used to save space for small lists.
	//     Sets can be encoded as intset or hashtable. The intset is a special encoding used for small sets composed solely of integers.
	//     Hashes can be encoded as zipmap or hashtable. The zipmap is a special encoding used for small hashes.
	//     Sorted Sets can be encoded as ziplist or skiplist format. As for the List type small sorted sets can be specially encoded using ziplist, while the skiplist encoding is the one that works with sorted sets of any size.
	// Different return values are used for different subcommands.
	//     Subcommands refcount and idletime return integers.
	//     Subcommand encoding returns a bulk reply.
	// If the object you try to inspect is missing, a null bulk reply is returned.
	// Object(subcommand string, arguments ...interface{}) (interface{}, error)

	// PERSIST key
	// Remove the expiration from a key
	// Integer reply, specifically:
	//     1 if the timeout was removed.
	//     0 if key does not exist or does not have an associated timeout.
	Persist(key interface{}) (int, error)

	// PEXPIRE key milliseconds
	// Set a key's time to live in milliseconds
	// Integer reply, specifically:
	//     1 if the timeout was set.
	//     0 if key does not exist or the timeout could not be set.
	PExpire(key string, milliseconds int) (int, error)

	// PEXPIREAT key milliseconds-timestamp
	// Set the expiration for a key as a UNIX timestamp specified in milliseconds
	// Integer reply, specifically:
	//     1 if the timeout was set.
	//     0 if key does not exist or the timeout could not be set (see: EXPIRE).
	PExpireAt(key string, millisecondTimestamp int) (int, error)

	// PTTL key
	// Get the time to live for a key in milliseconds
	// Integer reply: TTL in milliseconds, or a negative value in order to signal an error.
	Pttl(key string) (int, error)

	// RANDOMKEY
	// Return a random key from the keyspace
	// Bulk string reply: the random key, or nil when the database is empty.
	RandomKey() (interface{}, error)

	// RENAME key newkey
	// Rename a key
	// Simple string reply
	Rename(key, newkey string) (string, error)

	// RENAMENX key newkey
	// Rename a key, only if the new key does not exist
	// Integer reply, specifically:
	//     1 if key was renamed to newkey.
	//     0 if newkey already exists.
	RenameNx(key, newkey string) (int, error)

	// RESTORE key ttl serialized-value [REPLACE]
	// Create a key using the provided serialized value, previously obtained using DUMP.
	// Simple string reply: The command returns OK on success.
	Restore(key interface{}, ttl int, serializedValue interface{}, replace bool) (string, error)

	// SCAN cursor [MATCH pattern] [COUNT count]
	// Incrementally iterate the keys space

	// SORT key [BY pattern] [LIMIT offset count] [GET pattern [GET pattern ...]] [ASC|DESC] [ALPHA] [STORE destination]
	// Sort the elements in a list, set or sorted set
	// Array reply: list of sorted elements.

	// TTL key
	// Get the time to live for a key
	// Integer reply: TTL in seconds, or a negative value in order to signal an error.
	Ttl(key interface{}) (int, error)

	// TYPE key
	// Determine the type stored at key
	// Simple string reply: type of key, or none when key does not exist.
	Type(key interface{}) (string, error)

	// WAIT numslaves timeout
	// Wait for the synchronous replication of all the write commands sent in the context of the current connection
	// Integer reply: The command returns the number of slaves reached by all the writes performed in the context of the current connection.
	Wait(numslaves, timeout int) (int, error)
}
