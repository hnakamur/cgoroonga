package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func (c *Ctx) RecordAdd(table *Obj, key string) (recordID ID, added bool, err error) {
	var cKey *C.char
	var cKeyLen C.size_t
	if key != "" {
		cKey = C.CString(key)
		defer C.free(unsafe.Pointer(cKey))

		cKeyLen = C.strlen(cKey)
	}
	var cAdded C.int
	recordID = ID(C.grn_table_add(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(table)),
		unsafe.Pointer(cKey),
		C.uint(cKeyLen),
		&cAdded))
	if cAdded != 0 {
		added = true
	} else {
		added = false
	}
	return
}
