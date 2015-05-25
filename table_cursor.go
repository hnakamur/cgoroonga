package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"

type TableCursor struct {
	records *Records
	cCursor *C.grn_table_cursor
}

func (tc *TableCursor) Next() (ID, bool) {
	id := ID(C.grn_table_cursor_next(
		tc.records.db.context.cCtx,
		tc.cCursor))
	return id, id != ID(C.GRN_ID_NIL)
}

func (tc *TableCursor) Close() error {
	rc := C.grn_table_cursor_close(
		tc.records.db.context.cCtx,
		tc.cCursor)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}
