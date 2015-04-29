package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func (c *Ctx) TableCursorOpen(table *Obj, min, max string, offset, limit, flags int) (*TableCursor, error) {
	var cMin, cMax *C.char
	var cMinLen, cMaxLen C.size_t
	if min != "" {
		cMin = C.CString(min)
		defer C.free(unsafe.Pointer(cMin))
		cMinLen = C.strlen(cMin)
	}
	if max != "" {
		cMax = C.CString(max)
		defer C.free(unsafe.Pointer(cMax))
		cMaxLen = C.strlen(cMax)
	}

	cursor := C.grn_table_cursor_open(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(table)),
		unsafe.Pointer(cMin), C.uint(cMinLen),
		unsafe.Pointer(cMax), C.uint(cMaxLen),
		C.int(offset), C.int(limit), C.int(flags))
	if cursor == nil {
		return nil, errorFromRc(c.rc)
	}
	return (*TableCursor)(unsafe.Pointer(cursor)), nil
}

func (c *Ctx) TableCursorNext(tc *TableCursor) ID {
	return ID(C.grn_table_cursor_close(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.grn_table_cursor)(unsafe.Pointer(tc))))
}

func (c *Ctx) TableCursorClose(tc *TableCursor) error {
	rc := C.grn_table_cursor_close(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.grn_table_cursor)(unsafe.Pointer(tc)))
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}
