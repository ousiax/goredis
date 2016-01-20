// The MIT License (MIT)

// Copyright (c) 2016 Roy Xu

package redis

type respconnection interface {
	// AUTH password
	// Authenticate to the server
	// Simple string reply
	Auth(password string) (string, error)

	// ECHO message
	// Echo the given string
	// Simple string reply
	Echo(message string) (string, error)

	// PING
	// Ping the server
	// Simple string reply
	Ping() (string, error)

	// QUIT
	// Close the connection
	// Simple string reply: always OK.
	Quit() (string, error)

	// SELECT index
	// Change the selected database for the current connection
	// Simple string reply
	Select(index int) (string, error)
}
