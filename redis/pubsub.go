// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package redis

type respPubSub interface {
	// PSUBSCRIBE pattern [pattern ...] Listen for messages published to channels matching the given patterns
	// PUBSUB subcommand [argument [argument ...]] Inspect the state of the Pub/Sub subsystem
	// PUBLISH channel message Post a message to a channel
	// PUNSUBSCRIBE [pattern [pattern ...]] Stop listening for messages posted to channels matching the given patterns
	// SUBSCRIBE channel [channel ...] Listen for messages published to the given channels
	// UNSUBSCRIBE [channel [channel ...]] Stop listening for messages posted to the given channels
}
