package libubox

/*
#cgo LDFLAGS: -lubox
#include <libubox/uloop.h>

extern void uloopTimeoutHandlerProxy(struct uloop_timeout* t);
extern void uloopFDHandlerProxy(struct uloop_fd* u, unsigned int events);

// stub C function used as the handler for uloop timeout, invoke golang callback
// through uloopTimeoutHandlerProxy
void uloop_timeout_cb(struct uloop_timeout* t)
{
	uloopTimeoutHandlerProxy(t);
}

// stub C function used as the handler for uloop_fd, invoke golang callback
// through uloopFDHandlerProxy
void uloop_fd_cb(struct uloop_fd* u, unsigned int events)
{
	uloopFDHandlerProxy(u, events);
}
*/
import "C"
