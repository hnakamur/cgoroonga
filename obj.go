package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func unlinkObj(ctx *C.grn_ctx, obj *C.grn_obj) {
	C.grn_obj_unlink(ctx, obj)
}

func removeObj(ctx *C.grn_ctx, obj *C.grn_obj) error {
	rc := C.grn_obj_remove(ctx, obj)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func objPath(ctx *C.grn_ctx, obj *C.grn_obj) string {
	cPath := C.grn_obj_path(ctx, obj)
	if cPath == nil {
		return ""
	} else {
		return C.GoString(cPath)
	}
}

func objName(ctx *C.grn_ctx, obj *C.grn_obj) string {
	length := C.grn_obj_name(ctx, obj, nil, 0)
	if length == 0 {
		return ""
	}

	var buf *C.char
	buf = (*C.char)(C.malloc(C.size_t(unsafe.Sizeof(*buf)) * C.size_t(length)))
	defer C.free(unsafe.Pointer(buf))
	C.grn_obj_name(ctx, obj, buf, length)
	return C.GoString(buf)
}
