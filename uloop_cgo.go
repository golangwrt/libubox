package libubox

/*
#cgo LDFLAGS: -lubox
#include <stdlib.h>
#include <libubox/uloop.h>

extern void uloop_fd_cb(struct uloop_fd* u, unsigned int events);
extern void uloop_timeout_cb(struct uloop_timeout* t);
*/
import "C"
import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

type UloopFDFlag int

const (
	UloopRead        UloopFDFlag = C.ULOOP_READ
	UloopWrite       UloopFDFlag = C.ULOOP_WRITE
	UloopEdgeTrigger UloopFDFlag = C.ULOOP_EDGE_TRIGGER
	UloopBlocking    UloopFDFlag = C.ULOOP_BLOCKING
)

const (
	UFDReadBufferSize = 4096
)

// UloopFDHandler encapsulates uloop_fd_handler
//
// typedef void (*uloop_fd_handler)(struct uloop_fd *u, unsigned int events)
type UloopFDHandler func(ufd *UloopFD, events uint)

type UFDReadElem struct {
	Err  error
	Data []byte
	N    int
	ptr  *[]byte
}

// UloopFD encapsulates struct uloop_fd
type UloopFD struct {
	ptr *C.struct_uloop_fd
	h   UloopFDHandler
}

var (
	ufds struct {
		Elem map[*C.struct_uloop_fd]*UloopFD
		sync.RWMutex
	}
	ufdReadBufferPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, UFDReadBufferSize)
			return &b
		},
	}
)

func init() {
	ufds.Elem = make(map[*C.struct_uloop_fd]*UloopFD, 0)
}

func UloopInit() error {
	_, err := C.uloop_init()
	return err
}

func UloopRun() error {
	_, err := C.uloop_run()
	return err
}

// UloopEnd encapsualtes uloop_end
// set flag uloop_cancelled to true
func UloopEnd() error {
	_, err := C.uloop_end()
	return err
}

// Uloopdone encapsulates uloop_done
// do cleanup steps after uloop_run return
func UloopDone() error {
	_, err := C.uloop_done()
	return err
}

// NewUloopFD created an *UloopFD, and set the underlying
// UloopFDHandler to specified value
func NewUloopFD(fd int, h UloopFDHandler) *UloopFD {
	if h == nil {
		return nil
	}
	ufd := &UloopFD{
		h: h,
	}
	ufd.ptr = (*C.struct_uloop_fd)(C.calloc(1, C.sizeof_struct_uloop_fd))
	ufd.ptr.cb = C.uloop_fd_handler(C.uloop_fd_cb)
	ufd.ptr.fd = C.int(fd)
	return ufd
}

// Add encapsulates uloop_fd_add
//
// int uloop_fd_add(struct uloop_fd *sock, unsigned int flags)
func (ufd *UloopFD) Add(flag UloopFDFlag) error {
	ret, err := C.uloop_fd_add(ufd.ptr, C.uint(flag))
	if err != nil {
		return fmt.Errorf("ret: %d, %s", int(ret), err)
	}

	ufds.Lock()
	ufds.Elem[ufd.ptr] = ufd
	defer ufds.Unlock()

	return nil
}

// Delete encapsualtes uloop_fd_delete
//
// int uloop_fd_delete(struct uloop_fd *sock)
func (ufd *UloopFD) Delete() error {
	ret, err := C.uloop_fd_delete(ufd.ptr)
	if err != nil {
		return fmt.Errorf("ret: %d, %s", int(ret), err)
	}

	ufds.Lock()
	delete(ufds.Elem, ufd.ptr)
	ufds.Unlock()

	syscall.Close(int(ufd.ptr.fd))
	C.free(unsafe.Pointer(ufd.ptr))

	ufd.ptr = nil
	ufd.h = nil

	return nil
}

func (elem *UFDReadElem) Free() {
	if elem == nil || elem.ptr == nil {
		return
	}
	ufdReadBufferPool.Put(elem.ptr)
}

// Read read from the underlying file descriptor
// this call may block if there is non data ready.
// after finished using UFDReadElem, it's Free method must be called,
// otherwise memory leak will be occurred
func (ufd *UloopFD) Read() *UFDReadElem {
	buf := ufdReadBufferPool.Get().(*[]byte)
	n, err := syscall.Read(int(ufd.ptr.fd), *buf)
	return &UFDReadElem{
		Data: (*buf)[:n],
		Err:  err,
		N:    n,
		ptr:  buf,
	}
}

func (ufd *UloopFD) Write(data []byte) (n int, err error) {
	return syscall.Write(int(ufd.ptr.fd), data)
}

//export uloopFDHandlerProxy
func uloopFDHandlerProxy(ptr *C.struct_uloop_fd, events C.uint) {
	ufds.RLock()
	ufd, found := ufds.Elem[ptr]
	ufds.RUnlock()
	if !found {
		fmt.Fprintf(os.Stderr, "non ufd found for uloop_fd %ptr\n", ptr)
		return
	}
	ufd.h(ufd, uint(events))
}

//export uloopTimeoutHandlerProxy
func uloopTimeoutHandlerProxy(ptr *C.struct_uloop_timeout) {

}
