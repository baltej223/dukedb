// Package dukerror handles by this handler.
package dukerror

import (
	"context"
	"errors"
	"os"
	"strings"

	dukelog "github.com/baltej223/dukedb/log"
)

var (
	ErrNetwork         = errors.New("ErrNetwork: network error ")
	ErrTimeout         = errors.New("ErrTimeout: request timed out ")
	ErrWrongOwner      = errors.New("ErrWrongOwner: wrong owner ")
	ErrNodeDead        = errors.New("ErrNodeDead: node suspected dead ")
	ErrStorage         = errors.New("ErrStorage: storage error ")
	ErrKeyNotFound     = errors.New("ErrKeyNotFound: key not found ")
	ErrUnknownMessage  = errors.New("ErrUnknownMessage: unknown message type ")
	ErrReplication     = errors.New("ErrReplication: replication failed ")
	ErrAddressNotFree  = errors.New("ErrAddressNotFree: bind: bind address already in use ")
	ErrConnectionReset = errors.New("ErrConnectionReset: connection reset by peer ")
	ErrPermission      = errors.New("ErrPermission: operation not permitted ")

	ErrPeersDefinedForSeedNode = errors.New("ErrPeersDefinedForSeedNode: peers can't be defined for seed nodes ")
	ErrNonSeedNodeAndNoPeer    = errors.New("ErrNonSeedNodeAndNoPeer: atleast one peer (or the seed node) should be defined for a non seed node ")
	ErrClusterJoiningFailed    = errors.New("ErrClusterJoiningFailed: node: tried to join the cluster, but failed ")
)

func Normalize(err error) error {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return errors.Join(ErrTimeout, err)

	case errors.Is(err, os.ErrDeadlineExceeded):
		return errors.Join(ErrTimeout, err)

	case strings.Contains(err.Error(), "connection refused"):
		return errors.Join(ErrNetwork, err)

	case strings.Contains(err.Error(), "address already in use"):
		return errors.Join(ErrAddressNotFree, err)

	case strings.Contains(err.Error(), "connection reset by peer"):
		return errors.Join(ErrConnectionReset, err)

	case strings.Contains(err.Error(), "network is unreachable"):
		return errors.Join(ErrNetwork, err)

	case strings.Contains(err.Error(), "permission"):
		return errors.Join(ErrPermission, err)

	default:
		return err
	}
}

func Handle(err error) {
	if err == ErrPeersDefinedForSeedNode || err == ErrNonSeedNodeAndNoPeer {
		ShutdownGracefully(err)
	} else {
		dukelog.Error(err)
	}
}

func ShutdownGracefully(err error) {
	panic(err)
}
