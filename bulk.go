package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func BulkRewind(bulk *Obj) {
	C.go_grn_bulk_rewind((*C.struct__grn_obj)(unsafe.Pointer(bulk)))
}

func (c *Ctx) BulkReinit(bulk *Obj, size uint) error {
	rc := C.grn_bulk_reinit(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(bulk)),
		C.uint(size),
	)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (c *Ctx) BulkReserve(bulk *Obj, length uint) error {
	rc := C.grn_bulk_reserve(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(bulk)),
		C.uint(length),
	)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func BulkHead(bulk *Obj) (head string) {
	cHead := C.go_grn_bulk_head((*C.struct__grn_obj)(unsafe.Pointer(bulk)))
	cSize := C.go_grn_bulk_vsize((*C.struct__grn_obj)(unsafe.Pointer(bulk)))
	if cHead != nil {
		head = C.GoStringN(cHead, cSize)
	}
	return
}
