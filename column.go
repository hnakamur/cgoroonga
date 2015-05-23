package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

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
