package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func (c *Ctx) ObjClose(obj *Obj) error {
	rc := C.grn_obj_close(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(obj)),
	)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (c *Ctx) ObjSetValue(obj *Obj, recordID ID, value *Obj, flags int) error {
	rc := C.grn_obj_set_value(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(obj)),
		C.grn_id(recordID),
		(*C.struct__grn_obj)(unsafe.Pointer(value)),
		C.int(flags),
	)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (c *Ctx) ObjGetValue(obj *Obj, recordID ID, value *Obj) *Obj {
	return (*Obj)(unsafe.Pointer(C.grn_obj_get_value(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(obj)),
		C.grn_id(recordID),
		(*C.struct__grn_obj)(unsafe.Pointer(value)),
	)))
}

func (c *Ctx) ObjCloseDefer(err *error, obj *Obj) {
	err2 := c.ObjClose(obj)
	if err2 != nil && *err == nil {
		*err = err2
	}
}
