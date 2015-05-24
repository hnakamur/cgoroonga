package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import (
	"time"
	"unsafe"
)

type ID int

func (id ID) SetString(column *Column, str string) error {
	var cStr *C.char
	var cStrLen C.size_t
	if str != "" {
		cStr = C.CString(str)
		defer C.free(unsafe.Pointer(cStr))
		cStrLen = C.strlen(cStr)
	}

	cCtx := column.table.db.context.cCtx
	var buf C.grn_obj
	C.cgoroonga_text_init(&buf, 0)
	defer unlinkObj(cCtx, &buf)
	C.cgoroonga_text_put(cCtx, &buf, cStr, C.uint(cStrLen))
	rc := C.grn_obj_set_value(cCtx, column.cColumn, C.grn_id(id), &buf,
		C.GRN_OBJ_SET)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (id ID) GetString(column *Column) string {
	cCtx := column.table.db.context.cCtx
	var buf C.grn_obj
	C.cgoroonga_text_init(&buf, 0)
	defer unlinkObj(cCtx, &buf)
	C.grn_obj_get_value(cCtx, column.cColumn, C.grn_id(id), &buf)
	return C.GoStringN(C.cgoroonga_bulk_head(&buf),
		C.cgoroonga_bulk_vsize(&buf))
}

func (id ID) SetTime(column *Column, t time.Time) error {
	cCtx := column.table.db.context.cCtx
	var buf C.grn_obj
	C.cgoroonga_time_init(&buf, 0)
	defer unlinkObj(cCtx, &buf)
	// convert nano seconds to micro seconds
	usec := t.UnixNano() / 1000
	C.cgoroonga_time_set(cCtx, &buf, C.longlong(usec))
	rc := C.grn_obj_set_value(cCtx, column.cColumn, C.grn_id(id), &buf,
		C.GRN_OBJ_SET)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (id ID) GetTime(column *Column) time.Time {
	cCtx := column.table.db.context.cCtx
	var buf C.grn_obj
	C.cgoroonga_time_init(&buf, 0)
	defer unlinkObj(cCtx, &buf)
	C.grn_obj_get_value(cCtx, column.cColumn, C.grn_id(id), &buf)
	vsize := C.cgoroonga_bulk_vsize(&buf)
	if vsize == 0 {
		return time.Unix(0, 0)
	}
	usec := C.cgoroonga_int64_value(&buf)
	return time.Unix(int64(usec/1000000), int64((usec%1000000)*1000))
}
