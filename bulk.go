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

func BulkHead(bulk *Obj) (head string) {
	cHead := C.go_grn_bulk_head((*C.struct__grn_obj)(unsafe.Pointer(bulk)))
	cSize := C.go_grn_bulk_vsize((*C.struct__grn_obj)(unsafe.Pointer(bulk)))
	if cHead != nil {
		head = C.GoStringN(cHead, cSize)
	}
	return
}
