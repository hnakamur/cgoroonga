package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func CtxOpen(flags int) (*Ctx, error) {
	ctx := C.grn_ctx_open(C.int(flags))
	if ctx == nil {
		return nil, CtxOpenError
	}
	return (*Ctx)(unsafe.Pointer(ctx)), nil
}

func (c *Ctx) Close() error {
	rc := C.grn_ctx_close((*C.struct__grn_ctx)(unsafe.Pointer(c)))
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (c *Ctx) Get(name string) *Obj {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cNameLen := C.strlen(cName)
	obj := C.grn_ctx_get(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		cName, C.int(cNameLen))
	return (*Obj)(unsafe.Pointer(obj))
}

func (c *Ctx) At(id ID) *Obj {
	obj := C.grn_ctx_at(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		C.grn_id(id))
	return (*Obj)(unsafe.Pointer(obj))
}

func (c *Ctx) CloseDefer(err *error) {
	err2 := c.Close()
	if err2 != nil && *err == nil {
		*err = err2
	}
}
