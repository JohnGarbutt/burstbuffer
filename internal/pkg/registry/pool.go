package registry

import (
	"context"
	"encoding/json"
	"log"
)

type PoolRegistry interface {
	// Returns a summary of the current state of all pools, including the bricks in each pool
	Pools() ([]Pool, error)

	// TODO: Pool(name string) (Pool, error)

	// Update (or add) information on what bricks are present
	//
	// Note: it is possible to have bricks from multiple pools on a single host
	// If any bricks that were previously registered have gone away,
	// they will be removed, unless there is an associated BrickAllocation which will
	// cause the update to fail and returns an error.
	// If any bricks in the same pool have a different capacity,
	// the update fails and returns an error.
	UpdateHost(bricks []BrickInfo) error

	// While the process is still running this notifies others the host is up
	//
	// When a host is dead non of its bricks will get new volumes assigned,
	// and no bricks will get cleaned up until the next service start.
	// Error will be returned if the host info has not yet been written.
	KeepAliveHost(hostname string) error

	// Update a brick with allocation information.
	//
	// No update is made and an error is returned if:
	// any brick already has an allocation,
	// or any volume a brick is being assigned to already has an allocation,
	// or if any of the volumes do not exist
	// or if there is not exactly one primary brick.
	//
	// Note: you may assign multiple volumes in a single call, but all bricks
	// for a particular volume must be set in a single call
	AllocateBricksForVolume(volume Volume) ([]BrickAllocation, error)

	// Deallocate all bricks associated with the given volume
	//
	// No update is made and an error is returned if any of brick allocations don't match the current state.
	// If any host associated with one of the bricks is down, an error is returned and the deallocate is
	// recorded as requested and not executed.
	// Note: this returns as soon as deallocate is requested, doesn't wait for cleanup completion
	DeallocateBricks(volume VolumeName) error

	// This is called after DeallocateBricks has been processed
	HardDeleteAllocations(allocations []BrickAllocation) error

	// Get all the allocations for bricks associated with the specified hostname
	GetAllocationsForHost(hostname string) ([]BrickAllocation, error)

	// Get all the allocations for bricks associated with the specific volume
	GetAllocationsForVolume(volume VolumeName) ([]BrickAllocation, error)

	// Get information on a specific brick
	GetBrickInfo(hostname string, device string) (BrickInfo, error)

	// Returns a channel that reports all new brick allocations for given hostname
	//
	// The channel is closed when the context is cancelled or timeout.
	// Any errors in the watching log the issue and panic
	GetNewHostBrickAllocations(ctxt context.Context, hostname string) <-chan BrickAllocation
}

type Pool struct {
	// The pool is derived from all the reported bricks
	// It must only contain the characters A-Za-z0-9
	Name string // TODO: should we create PoolName type?

	// Returns all unallocated bricks in this pool associated with a live host
	AvailableBricks []BrickInfo

	// Returns all brick allocations for this pool
	AllocatedBricks []BrickAllocation

	// This is the allocation unit for the pool
	// It is the minimum size of any registered brick
	GranularityGB uint

	// List of all hosts that report bricks in this pool
	Hosts map[string]HostInfo
}

func (pool Pool) String() string {
	poolString, err := json.Marshal(pool)
	if err != nil {
		log.Fatal(err)
	}
	return string(poolString)
}

type HostInfo struct {
	// It must only contain the characters "A-Za-z0-9."
	Hostname string

	// True if data accelerator process is thought to be running
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

type BrickAllocation struct {
	// Bricks are identified by device and hostname
	// It must only contain the characters A-Za-z0-9
	Device string

	// It must only contain the characters "A-Za-z0-9."
	Hostname string

	// Name of the volume that owns the brick
	AllocatedVolume VolumeName

	// 0 index allocation is the primary brick,
	// which is responsible for provisioning the associated volume
	AllocatedIndex uint

	// If any allocation sent to deallocate has a host that isn't
	// alive, this flag is set rather than have allocations removed.
	// A host should check for any allocations
	DeallocateRequested bool
}
