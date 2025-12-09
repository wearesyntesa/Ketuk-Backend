package scheduler

import "sync/atomic"

var unblockState atomic.Bool

func init() {
	unblockState.Store(false)
}

func IsUnblockEnabled() bool {
	return unblockState.Load()
}

func EnableUnblock() {
	unblockState.Store(true)
}

func DisableUnblock() {
	unblockState.Store(false)
}
