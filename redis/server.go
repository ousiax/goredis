// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package redis

type respServer interface {
	// BGREWRITEAOF Asynchronously rewrite the append-only file
	// BGSAVE Asynchronously save the dataset to disk
	// CLIENT KILL [ip:port] [ID client-id] [TYPE normal|slave|pubsub] [ADDR ip:port] [SKIPME yes/no] Kill the connection of a client
	// CLIENT LIST Get the list of client connections
	// CLIENT GETNAME Get the current connection name
	// CLIENT PAUSE timeout Stop processing commands from clients for some time
	// CLIENT SETNAME connection-name Set the current connection name
	// COMMAND Get array of Redis command details
	// COMMAND COUNT Get total number of Redis commands
	// COMMAND GETKEYS Extract keys given a full Redis command
	// COMMAND INFO command-name [command-name ...] Get array of specific Redis command details
	// CONFIG GET parameter Get the value of a configuration parameter
	// CONFIG REWRITE Rewrite the configuration file with the in memory configuration
	// CONFIG SET parameter value Set a configuration parameter to the given value
	// CONFIG RESETSTAT Reset the stats returned by INFO
	// DBSIZE Return the number of keys in the selected database
	// DEBUG OBJECT key Get debugging information about a key
	// DEBUG SEGFAULT Make the server crash
	// FLUSHALL Remove all keys from all databases
	// FLUSHDB Remove all keys from the current database
	// INFO [section] Get information and statistics about the server
	// LASTSAVE Get the UNIX time stamp of the last successful save to disk
	// MONITOR Listen for all requests received by the server in real time
	// ROLE Return the role of the instance in the context of replication
	// SAVE Synchronously save the dataset to disk
	// SHUTDOWN [NOSAVE] [SAVE] Synchronously save the dataset to disk and then shut down the server
	// SLAVEOF host port Make the server a slave of another instance, or promote it as master
	// SLOWLOG subcommand [argument] Manages the Redis slow queries log
	// SYNC Internal command used for replication
	// TIME Return the current server time
}
