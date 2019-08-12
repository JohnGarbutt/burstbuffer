package model

type BrickHostName string

type BrickHostInfo struct {
	Name BrickHostName

	// Returns all bricks
	Bricks []BrickInfo

	// True if allowing new volumes to use bricks from this host
	Enabled bool
}

type BrickHostStatus struct {
	Info BrickHostInfo

	// True is current keepalive key exists
	Alive bool
}

type BrickInfo struct {
	// Bricks are identified by device and hostname
	// It must only contain the characters A-Za-z0-9
	Device string

	// It must only contain the characters "A-Za-z0-9."
	Hostname string

	// The bool a brick is associated with
	// It must only contain the characters A-Za-z0-9
	PoolName string

	// Size of the brick, defines the pool granularity
	CapacityGB uint
}
