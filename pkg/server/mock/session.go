// Code generated by moq; DO NOT EDIT.
// gitlab.diskarte.net/engineering/moqc fork from github.com/matryer/moq

package mock

import (
	"net"
	"sync"

	"github.com/DarthPestilane/easytcp"
)

var _ easytcp.Session = &Session{}

// Session is a mock implementation of server.Session.
type Session struct {
	// AfterCloseHookFn mocks the AfterCloseHook method.
	AfterCloseHookFn func() <-chan struct{}

	// AfterCreateHookFn mocks the AfterCreateHook method.
	AfterCreateHookFn func() <-chan struct{}

	// AllocateContextFn mocks the AllocateContext method.
	AllocateContextFn func() easytcp.Context

	// CloseFn mocks the Close method.
	CloseFn func()

	// CodecFn mocks the Codec method.
	CodecFn func() easytcp.Codec

	// ConnFn mocks the Conn method.
	ConnFn func() net.Conn

	// IDFn mocks the ID method.
	IDFn func() interface{}

	// SendFn mocks the Send method.
	SendFn func(ctx easytcp.Context) bool

	// SetIDFn mocks the SetID method.
	SetIDFn func(id interface{})

	// calls tracks calls to the methods.
	calls struct {
		// AfterCloseHook holds details about calls to the AfterCloseHook method.
		AfterCloseHook []struct {
		}
		// AfterCreateHook holds details about calls to the AfterCreateHook method.
		AfterCreateHook []struct {
		}
		// AllocateContext holds details about calls to the AllocateContext method.
		AllocateContext []struct {
		}
		// Close holds details about calls to the Close method.
		Close []struct {
		}
		// Codec holds details about calls to the Codec method.
		Codec []struct {
		}
		// Conn holds details about calls to the Conn method.
		Conn []struct {
		}
		// ID holds details about calls to the ID method.
		ID []struct {
		}
		// Send holds details about calls to the Send method.
		Send []struct {
			Ctx easytcp.Context
		}
		// SetID holds details about calls to the SetID method.
		SetID []struct {
			ID interface{}
		}
	}
	lockAfterCloseHook  sync.RWMutex
	lockAfterCreateHook sync.RWMutex
	lockAllocateContext sync.RWMutex
	lockClose           sync.RWMutex
	lockCodec           sync.RWMutex
	lockConn            sync.RWMutex
	lockID              sync.RWMutex
	lockSend            sync.RWMutex
	lockSetID           sync.RWMutex
}

// AfterCloseHook calls AfterCloseHookFn.
func (mock *Session) AfterCloseHook() <-chan struct{} {
	callInfo := struct {
	}{}
	mock.lockAfterCloseHook.Lock()
	mock.calls.AfterCloseHook = append(mock.calls.AfterCloseHook, callInfo)
	mock.lockAfterCloseHook.Unlock()
	if mock.AfterCloseHookFn == nil {
		var (
			out1 <-chan struct{}
		)
		return out1
	}
	return mock.AfterCloseHookFn()
}

// AfterCloseHookCalls gets all the calls that were made to AfterCloseHook.
// Check the length with:
//
//	len(mockedSession.AfterCloseHookCalls())
func (mock *Session) AfterCloseHookCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockAfterCloseHook.RLock()
	calls = mock.calls.AfterCloseHook
	mock.lockAfterCloseHook.RUnlock()
	return calls
}

// AfterCreateHook calls AfterCreateHookFn.
func (mock *Session) AfterCreateHook() <-chan struct{} {
	callInfo := struct {
	}{}
	mock.lockAfterCreateHook.Lock()
	mock.calls.AfterCreateHook = append(mock.calls.AfterCreateHook, callInfo)
	mock.lockAfterCreateHook.Unlock()
	if mock.AfterCreateHookFn == nil {
		var (
			out1 <-chan struct{}
		)
		return out1
	}
	return mock.AfterCreateHookFn()
}

// AfterCreateHookCalls gets all the calls that were made to AfterCreateHook.
// Check the length with:
//
//	len(mockedSession.AfterCreateHookCalls())
func (mock *Session) AfterCreateHookCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockAfterCreateHook.RLock()
	calls = mock.calls.AfterCreateHook
	mock.lockAfterCreateHook.RUnlock()
	return calls
}

// AllocateContext calls AllocateContextFn.
func (mock *Session) AllocateContext() easytcp.Context {
	callInfo := struct {
	}{}
	mock.lockAllocateContext.Lock()
	mock.calls.AllocateContext = append(mock.calls.AllocateContext, callInfo)
	mock.lockAllocateContext.Unlock()
	if mock.AllocateContextFn == nil {
		var (
			out1 easytcp.Context
		)
		return out1
	}
	return mock.AllocateContextFn()
}

// AllocateContextCalls gets all the calls that were made to AllocateContext.
// Check the length with:
//
//	len(mockedSession.AllocateContextCalls())
func (mock *Session) AllocateContextCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockAllocateContext.RLock()
	calls = mock.calls.AllocateContext
	mock.lockAllocateContext.RUnlock()
	return calls
}

// Close calls CloseFn.
func (mock *Session) Close() {
	callInfo := struct {
	}{}
	mock.lockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	mock.lockClose.Unlock()
	if mock.CloseFn == nil {
		return
	}
	mock.CloseFn()
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//
//	len(mockedSession.CloseCalls())
func (mock *Session) CloseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockClose.RLock()
	calls = mock.calls.Close
	mock.lockClose.RUnlock()
	return calls
}

// Codec calls CodecFn.
func (mock *Session) Codec() easytcp.Codec {
	callInfo := struct {
	}{}
	mock.lockCodec.Lock()
	mock.calls.Codec = append(mock.calls.Codec, callInfo)
	mock.lockCodec.Unlock()
	if mock.CodecFn == nil {
		var (
			out1 easytcp.Codec
		)
		return out1
	}
	return mock.CodecFn()
}

// CodecCalls gets all the calls that were made to Codec.
// Check the length with:
//
//	len(mockedSession.CodecCalls())
func (mock *Session) CodecCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockCodec.RLock()
	calls = mock.calls.Codec
	mock.lockCodec.RUnlock()
	return calls
}

// Conn calls ConnFn.
func (mock *Session) Conn() net.Conn {
	callInfo := struct {
	}{}
	mock.lockConn.Lock()
	mock.calls.Conn = append(mock.calls.Conn, callInfo)
	mock.lockConn.Unlock()
	if mock.ConnFn == nil {
		var (
			out1 net.Conn
		)
		return out1
	}
	return mock.ConnFn()
}

// ConnCalls gets all the calls that were made to Conn.
// Check the length with:
//
//	len(mockedSession.ConnCalls())
func (mock *Session) ConnCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockConn.RLock()
	calls = mock.calls.Conn
	mock.lockConn.RUnlock()
	return calls
}

// ID calls IDFn.
func (mock *Session) ID() interface{} {
	callInfo := struct {
	}{}
	mock.lockID.Lock()
	mock.calls.ID = append(mock.calls.ID, callInfo)
	mock.lockID.Unlock()
	if mock.IDFn == nil {
		var (
			out1 interface{}
		)
		return out1
	}
	return mock.IDFn()
}

// IDCalls gets all the calls that were made to ID.
// Check the length with:
//
//	len(mockedSession.IDCalls())
func (mock *Session) IDCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockID.RLock()
	calls = mock.calls.ID
	mock.lockID.RUnlock()
	return calls
}

// Send calls SendFn.
func (mock *Session) Send(ctx easytcp.Context) bool {
	callInfo := struct {
		Ctx easytcp.Context
	}{
		Ctx: ctx,
	}
	mock.lockSend.Lock()
	mock.calls.Send = append(mock.calls.Send, callInfo)
	mock.lockSend.Unlock()
	if mock.SendFn == nil {
		var (
			out1 bool
		)
		return out1
	}
	return mock.SendFn(ctx)
}

// SendCalls gets all the calls that were made to Send.
// Check the length with:
//
//	len(mockedSession.SendCalls())
func (mock *Session) SendCalls() []struct {
	Ctx easytcp.Context
} {
	var calls []struct {
		Ctx easytcp.Context
	}
	mock.lockSend.RLock()
	calls = mock.calls.Send
	mock.lockSend.RUnlock()
	return calls
}

// SetID calls SetIDFn.
func (mock *Session) SetID(id interface{}) {
	callInfo := struct {
		ID interface{}
	}{
		ID: id,
	}
	mock.lockSetID.Lock()
	mock.calls.SetID = append(mock.calls.SetID, callInfo)
	mock.lockSetID.Unlock()
	if mock.SetIDFn == nil {
		return
	}
	mock.SetIDFn(id)
}

// SetIDCalls gets all the calls that were made to SetID.
// Check the length with:
//
//	len(mockedSession.SetIDCalls())
func (mock *Session) SetIDCalls() []struct {
	ID interface{}
} {
	var calls []struct {
		ID interface{}
	}
	mock.lockSetID.RLock()
	calls = mock.calls.SetID
	mock.lockSetID.RUnlock()
	return calls
}
