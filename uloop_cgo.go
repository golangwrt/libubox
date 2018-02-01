package libubox

/*
#cgo LDFLAGS: -lubox
#include <libubox/uloop.h>

*/
import "C"

const (
	UloopRead = C.ULOOP_READ
	UloopWrite = C.ULOOP_WRITE
	UloopEdgeTrigger = C.ULOOP_EDGE_TRIGGER
	UloopBlocking = C.ULOOP_BLOCKING
)

func UloopInit() error {
	_, err := C.uloop_init()
	return err
}

func UloopRun() error {
	_, err := C.uloop_run()
	return err
}


