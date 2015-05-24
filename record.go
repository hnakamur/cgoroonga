package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

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
