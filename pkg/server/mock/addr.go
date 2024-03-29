// Code generated by moq; DO NOT EDIT.
// gitlab.diskarte.net/engineering/moqc fork from github.com/matryer/moq

package mock

import (
	"net"
	"sync"
)

var _ net.Addr = &Addr{}

// Addr is a mock implementation of server.Addr.
type Addr struct {
	// NetworkFn mocks the Network method.
	NetworkFn func() string

	// StringFn mocks the String method.
	StringFn func() string

	// calls tracks calls to the methods.
	calls struct {
		// Network holds details about calls to the Network method.
		Network []struct {
		}
		// String holds details about calls to the String method.
		String []struct {
		}
	}
	lockNetwork sync.RWMutex
	lockString  sync.RWMutex
}

// Network calls NetworkFn.
func (mock *Addr) Network() string {
	callInfo := struct {
	}{}
	mock.lockNetwork.Lock()
	mock.calls.Network = append(mock.calls.Network, callInfo)
	mock.lockNetwork.Unlock()
	if mock.NetworkFn == nil {
		var (
			out1 string
		)
		return out1
	}
	return mock.NetworkFn()
}

// NetworkCalls gets all the calls that were made to Network.
// Check the length with:
//
//	len(mockedAddr.NetworkCalls())
func (mock *Addr) NetworkCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockNetwork.RLock()
	calls = mock.calls.Network
	mock.lockNetwork.RUnlock()
	return calls
}

// String calls StringFn.
func (mock *Addr) String() string {
	callInfo := struct {
	}{}
	mock.lockString.Lock()
	mock.calls.String = append(mock.calls.String, callInfo)
	mock.lockString.Unlock()
	if mock.StringFn == nil {
		var (
			out1 string
		)
		return out1
	}
	return mock.StringFn()
}

// StringCalls gets all the calls that were made to String.
// Check the length with:
//
//	len(mockedAddr.StringCalls())
func (mock *Addr) StringCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockString.RLock()
	calls = mock.calls.String
	mock.lockString.RUnlock()
	return calls
}
