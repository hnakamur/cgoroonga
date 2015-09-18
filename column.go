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

type Column struct {
	table   *Table
	cColumn *C.grn_obj
}

func (c *Column) Name() string {
	cCtx := c.table.db.context.cCtx
	length := C.grn_column_name(cCtx, c.cColumn, nil, 0)
	if length == 0 {
		return ""
	}

	var buf *C.char
	buf = (*C.char)(C.malloc(C.size_t(unsafe.Sizeof(*buf)) * C.size_t(length)))
	defer C.free(unsafe.Pointer(buf))
	C.grn_column_name(cCtx, c.cColumn, buf, length)
	return C.GoString(buf)
}

func (c *Column) Path() string {
	return objPath(c.table.db.context.cCtx, c.cColumn)
}

func (c *Column) Close() {
	if c.cColumn == nil {
		return
	}
	unlinkObj(c.table.db.context.cCtx, c.cColumn)
	c.cColumn = nil
}

func (c *Column) Remove() error {
	if c.cColumn == nil {
		return InvalidArgumentError
	}
	err := removeObj(c.table.db.context.cCtx, c.cColumn)
	c.cColumn = nil
	return err
}

func (c *Column) SetString(id ID, str string) error {
	var cStr *C.char
	var cStrLen C.size_t
	if str != "" {
		cStr = C.CString(str)
		defer C.free(unsafe.Pointer(cStr))
		cStrLen = C.strlen(cStr)
	}

	cCtx := c.table.db.context.cCtx
	var buf C.grn_obj
	C.cgoroonga_text_init(&buf, 0)
	defer unlinkObj(cCtx, &buf)
	C.cgoroonga_text_put(cCtx, &buf, cStr, C.uint(cStrLen))
	rc := C.grn_obj_set_value(cCtx, c.cColumn, C.grn_id(id), &buf,
		C.GRN_OBJ_SET)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (c *Column) GetString(id ID) string {
	cCtx := c.table.db.context.cCtx
	var buf C.grn_obj
	C.cgoroonga_text_init(&buf, 0)
	defer unlinkObj(cCtx, &buf)
	C.grn_obj_get_value(cCtx, c.cColumn, C.grn_id(id), &buf)
	return C.GoStringN(C.cgoroonga_bulk_head(&buf),
		C.cgoroonga_bulk_vsize(&buf))
}

func (c *Column) SetTime(id ID, t time.Time) error {
	cCtx := c.table.db.context.cCtx
	var buf C.grn_obj
	C.cgoroonga_time_init(&buf, 0)
	defer unlinkObj(cCtx, &buf)
	// convert nano seconds to micro seconds
	usec := t.UnixNano() / 1000
	C.cgoroonga_time_set(cCtx, &buf, C.longlong(usec))
	rc := C.grn_obj_set_value(cCtx, c.cColumn, C.grn_id(id), &buf,
		C.GRN_OBJ_SET)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (c *Column) GetTime(id ID) time.Time {
	cCtx := c.table.db.context.cCtx
	var buf C.grn_obj
	C.cgoroonga_time_init(&buf, 0)
	defer unlinkObj(cCtx, &buf)
	C.grn_obj_get_value(cCtx, c.cColumn, C.grn_id(id), &buf)
	vsize := C.cgoroonga_bulk_vsize(&buf)
	if vsize == 0 {
		return time.Unix(0, 0)
	}
	usec := C.cgoroonga_int64_value(&buf)
	return time.Unix(int64(usec/1000000), int64((usec%1000000)*1000))
}
