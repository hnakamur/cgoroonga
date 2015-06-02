package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

type Table struct {
	*Records
}

func (t *Table) CreateColumn(name, path string, flags, columnType int) (*Column, error) {
	var cName *C.char
	var cNameLen C.size_t
	if name != "" {
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
		cNameLen = C.strlen(cName)
	}

	var cPath *C.char
	if path != "" {
		cPath = C.CString(path)
		defer C.free(unsafe.Pointer(cPath))
	}

	cCtx := t.db.context.cCtx
	columnTypeObj := C.grn_ctx_at(cCtx, C.grn_id(columnType))
	if columnTypeObj == nil {
		return nil, InvalidArgumentError
	}
	cColumn := C.grn_column_create(cCtx, t.cRecords, cName, C.uint(cNameLen),
		cPath, C.grn_obj_flags(flags), columnTypeObj)
	if cColumn == nil {
		return nil, errorFromRc(cCtx.rc)
	}
	column := &Column{table: t, cColumn: cColumn}
	t.addColumnToMap(name, column)
	return column, nil
}

func (t *Table) OpenOrCreateColumn(name, path string, flags, columnType int) (*Column, error) {
	column, err := t.OpenColumn(name)
	if err != nil {
		if err != NotFoundError {
			return nil, err
		}
		column, err = t.CreateColumn(name, path, flags, columnType)
	}
	return column, err
}

func (t *Table) IsLocked() bool {
	locked := C.grn_obj_is_locked(t.db.context.cCtx, t.cRecords)
	return locked != 0
}

func (t *Table) Lock(seconds int) error {
	rc := C.grn_obj_lock(t.db.context.cCtx, t.cRecords, C.GRN_ID_NIL, C.int(seconds))
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (t *Table) Unlock() error {
	rc := C.grn_obj_unlock(t.db.context.cCtx, t.cRecords, C.GRN_ID_NIL)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}
