// Code generated by moq; DO NOT EDIT.
// gitlab.diskarte.net/engineering/moqc fork from github.com/matryer/moq

package mock

import (
	"net"
	"sync"
	"time"
)

var _ net.Conn = &Conn{}

// Conn is a mock implementation of server.Conn.
type Conn struct {
	// CloseFn mocks the Close method.
	CloseFn func() error

	// LocalAddrFn mocks the LocalAddr method.
	LocalAddrFn func() net.Addr

	// ReadFn mocks the Read method.
	ReadFn func(b []byte) (int, error)

	// RemoteAddrFn mocks the RemoteAddr method.
	RemoteAddrFn func() net.Addr

	// SetDeadlineFn mocks the SetDeadline method.
	SetDeadlineFn func(t time.Time) error

	// SetReadDeadlineFn mocks the SetReadDeadline method.
	SetReadDeadlineFn func(t time.Time) error

	// SetWriteDeadlineFn mocks the SetWriteDeadline method.
	SetWriteDeadlineFn func(t time.Time) error

	// WriteFn mocks the Write method.
	WriteFn func(b []byte) (int, error)

	// calls tracks calls to the methods.
	calls struct {
		// Close holds details about calls to the Close method.
		Close []struct {
		}
		// LocalAddr holds details about calls to the LocalAddr method.
		LocalAddr []struct {
		}
		// Read holds details about calls to the Read method.
		Read []struct {
			B []byte
		}
		// RemoteAddr holds details about calls to the RemoteAddr method.
		RemoteAddr []struct {
		}
		// SetDeadline holds details about calls to the SetDeadline method.
		SetDeadline []struct {
			T time.Time
		}
		// SetReadDeadline holds details about calls to the SetReadDeadline method.
		SetReadDeadline []struct {
			T time.Time
		}
		// SetWriteDeadline holds details about calls to the SetWriteDeadline method.
		SetWriteDeadline []struct {
			T time.Time
		}
		// Write holds details about calls to the Write method.
		Write []struct {
			B []byte
		}
	}
	lockClose            sync.RWMutex
	lockLocalAddr        sync.RWMutex
	lockRead             sync.RWMutex
	lockRemoteAddr       sync.RWMutex
	lockSetDeadline      sync.RWMutex
	lockSetReadDeadline  sync.RWMutex
	lockSetWriteDeadline sync.RWMutex
	lockWrite            sync.RWMutex
}

// Close calls CloseFn.
func (mock *Conn) Close() error {
	callInfo := struct {
	}{}
	mock.lockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	mock.lockClose.Unlock()
	if mock.CloseFn == nil {
		var (
			out1 error
		)
		return out1
	}
	return mock.CloseFn()
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//
//	len(mockedConn.CloseCalls())
func (mock *Conn) CloseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockClose.RLock()
	calls = mock.calls.Close
	mock.lockClose.RUnlock()
	return calls
}

// LocalAddr calls LocalAddrFn.
func (mock *Conn) LocalAddr() net.Addr {
	callInfo := struct {
	}{}
	mock.lockLocalAddr.Lock()
	mock.calls.LocalAddr = append(mock.calls.LocalAddr, callInfo)
	mock.lockLocalAddr.Unlock()
	if mock.LocalAddrFn == nil {
		var (
			out1 net.Addr
		)
		return out1
	}
	return mock.LocalAddrFn()
}

// LocalAddrCalls gets all the calls that were made to LocalAddr.
// Check the length with:
//
//	len(mockedConn.LocalAddrCalls())
func (mock *Conn) LocalAddrCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockLocalAddr.RLock()
	calls = mock.calls.LocalAddr
	mock.lockLocalAddr.RUnlock()
	return calls
}

// Read calls ReadFn.
func (mock *Conn) Read(b []byte) (int, error) {
	callInfo := struct {
		B []byte
	}{
		B: b,
	}
	mock.lockRead.Lock()
	mock.calls.Read = append(mock.calls.Read, callInfo)
	mock.lockRead.Unlock()
	if mock.ReadFn == nil {
		var (
			n   int
			err error
		)
		return n, err
	}
	return mock.ReadFn(b)
}

// ReadCalls gets all the calls that were made to Read.
// Check the length with:
//
//	len(mockedConn.ReadCalls())
func (mock *Conn) ReadCalls() []struct {
	B []byte
} {
	var calls []struct {
		B []byte
	}
	mock.lockRead.RLock()
	calls = mock.calls.Read
	mock.lockRead.RUnlock()
	return calls
}

// RemoteAddr calls RemoteAddrFn.
func (mock *Conn) RemoteAddr() net.Addr {
	callInfo := struct {
	}{}
	mock.lockRemoteAddr.Lock()
	mock.calls.RemoteAddr = append(mock.calls.RemoteAddr, callInfo)
	mock.lockRemoteAddr.Unlock()
	if mock.RemoteAddrFn == nil {
		var (
			out1 net.Addr
		)
		return out1
	}
	return mock.RemoteAddrFn()
}

// RemoteAddrCalls gets all the calls that were made to RemoteAddr.
// Check the length with:
//
//	len(mockedConn.RemoteAddrCalls())
func (mock *Conn) RemoteAddrCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockRemoteAddr.RLock()
	calls = mock.calls.RemoteAddr
	mock.lockRemoteAddr.RUnlock()
	return calls
}

// SetDeadline calls SetDeadlineFn.
func (mock *Conn) SetDeadline(t time.Time) error {
	callInfo := struct {
		T time.Time
	}{
		T: t,
	}
	mock.lockSetDeadline.Lock()
	mock.calls.SetDeadline = append(mock.calls.SetDeadline, callInfo)
	mock.lockSetDeadline.Unlock()
	if mock.SetDeadlineFn == nil {
		var (
			out1 error
		)
		return out1
	}
	return mock.SetDeadlineFn(t)
}

// SetDeadlineCalls gets all the calls that were made to SetDeadline.
// Check the length with:
//
//	len(mockedConn.SetDeadlineCalls())
func (mock *Conn) SetDeadlineCalls() []struct {
	T time.Time
} {
	var calls []struct {
		T time.Time
	}
	mock.lockSetDeadline.RLock()
	calls = mock.calls.SetDeadline
	mock.lockSetDeadline.RUnlock()
	return calls
}

// SetReadDeadline calls SetReadDeadlineFn.
func (mock *Conn) SetReadDeadline(t time.Time) error {
	callInfo := struct {
		T time.Time
	}{
		T: t,
	}
	mock.lockSetReadDeadline.Lock()
	mock.calls.SetReadDeadline = append(mock.calls.SetReadDeadline, callInfo)
	mock.lockSetReadDeadline.Unlock()
	if mock.SetReadDeadlineFn == nil {
		var (
			out1 error
		)
		return out1
	}
	return mock.SetReadDeadlineFn(t)
}

// SetReadDeadlineCalls gets all the calls that were made to SetReadDeadline.
// Check the length with:
//
//	len(mockedConn.SetReadDeadlineCalls())
func (mock *Conn) SetReadDeadlineCalls() []struct {
	T time.Time
} {
	var calls []struct {
		T time.Time
	}
	mock.lockSetReadDeadline.RLock()
	calls = mock.calls.SetReadDeadline
	mock.lockSetReadDeadline.RUnlock()
	return calls
}

// SetWriteDeadline calls SetWriteDeadlineFn.
func (mock *Conn) SetWriteDeadline(t time.Time) error {
	callInfo := struct {
		T time.Time
	}{
		T: t,
	}
	mock.lockSetWriteDeadline.Lock()
	mock.calls.SetWriteDeadline = append(mock.calls.SetWriteDeadline, callInfo)
	mock.lockSetWriteDeadline.Unlock()
	if mock.SetWriteDeadlineFn == nil {
		var (
			out1 error
		)
		return out1
	}
	return mock.SetWriteDeadlineFn(t)
}

// SetWriteDeadlineCalls gets all the calls that were made to SetWriteDeadline.
// Check the length with:
//
//	len(mockedConn.SetWriteDeadlineCalls())
func (mock *Conn) SetWriteDeadlineCalls() []struct {
	T time.Time
} {
	var calls []struct {
		T time.Time
	}
	mock.lockSetWriteDeadline.RLock()
	calls = mock.calls.SetWriteDeadline
	mock.lockSetWriteDeadline.RUnlock()
	return calls
}

// Write calls WriteFn.
func (mock *Conn) Write(b []byte) (int, error) {
	callInfo := struct {
		B []byte
	}{
		B: b,
	}
	mock.lockWrite.Lock()
	mock.calls.Write = append(mock.calls.Write, callInfo)
	mock.lockWrite.Unlock()
	if mock.WriteFn == nil {
		var (
			n   int
			err error
		)
		return n, err
	}
	return mock.WriteFn(b)
}

// WriteCalls gets all the calls that were made to Write.
// Check the length with:
//
//	len(mockedConn.WriteCalls())
func (mock *Conn) WriteCalls() []struct {
	B []byte
} {
	var calls []struct {
		B []byte
	}
	mock.lockWrite.RLock()
	calls = mock.calls.Write
	mock.lockWrite.RUnlock()
	return calls
}
