package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func (c *Ctx) ObjUnlink(obj *Obj) error {
	C.grn_obj_unlink(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(obj)),
	)
	if c.rc != SUCCESS {
		return errorFromRc(c.rc)
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

func (c *Ctx) ObjUnlinkDefer(err *error, obj *Obj) {
	err2 := c.ObjUnlink(obj)
	if err2 != nil && *err == nil {
		*err = err2
	}
}

func (c *Ctx) ObjColumn(table *Obj, name string) (*Obj, error) {
	var cName *C.char
	var cNameLen C.size_t
	if name != "" {
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
		cNameLen = C.strlen(cName)
	}

	column := C.grn_obj_column(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(table)),
		cName, C.uint(cNameLen))
	if column == nil {
		if c.rc != SUCCESS {
			return nil, errorFromRc(c.rc)
		} else {
			return nil, ObjColumnError
		}
	}
	return (*Obj)(unsafe.Pointer(column)), nil
}
