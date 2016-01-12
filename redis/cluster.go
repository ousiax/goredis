// The MIT License (MIT)
//
// Copyright (c) 2016 Roy Xu

package redis

type respCluster interface {
	// CLUSTER ADDSLOTS slot [slot ...] Assign new hash slots to receiving node
	// CLUSTER COUNT-FAILURE-REPORTS node-id Return the number of failure reports active for a given node
	// CLUSTER COUNTKEYSINSLOT slot Return the number of local keys in the specified hash slot
	// CLUSTER DELSLOTS slot [slot ...] Set hash slots as unbound in receiving node
	// CLUSTER FAILOVER [FORCE|TAKEOVER] Forces a slave to perform a manual failover of its master.
	// CLUSTER FORGET node-id Remove a node from the nodes table
	// CLUSTER GETKEYSINSLOT slot count Return local key names in the specified hash slot
	// CLUSTER INFO Provides info about Redis Cluster node state
	// CLUSTER KEYSLOT key Returns the hash slot of the specified key
	// CLUSTER MEET ip port Force a node cluster to handshake with another node
	// CLUSTER NODES Get Cluster config for the node
	// CLUSTER REPLICATE node-id Reconfigure a node as a slave of the specified master node
	// CLUSTER RESET [HARD|SOFT] Reset a Redis Cluster node
	// CLUSTER SAVECONFIG Forces the node to save cluster state on disk
	// CLUSTER SET-CONFIG-EPOCH config-epoch Set the configuration epoch in a new node
	// CLUSTER SETSLOT slot IMPORTING|MIGRATING|STABLE|NODE [node-id] Bind an hash slot to a specific node
	// CLUSTER SLAVES node-id List slave nodes of the specified master node
	// CLUSTER SLOTS Get array of Cluster slot to node mappings
	// READONLY Enables read queries for a connection to a cluster slave node
	// READWRITE Disables read queries for a connection to a cluster slave node
}
