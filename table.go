package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

type Table struct {
	db      *DB
	cTable  *C.grn_obj
	columns map[string]*Column
}

func (t *Table) Path() string {
	return objPath(t.db.context.cCtx, t.cTable)
}

func (t *Table) Name() string {
	return objName(t.db.context.cCtx, t.cTable)
}

func (t *Table) Close() {
	for name, column := range t.columns {
		column.Close()
		delete(t.columns, name)
	}

	if t.cTable == nil {
		return
	}
	unlinkObj(t.db.context.cCtx, t.cTable)
	t.cTable = nil
}

func (t *Table) Remove() error {
	if t.cTable == nil {
		return InvalidArgumentError
	}
	err := removeObj(t.db.context.cCtx, t.cTable)
	t.cTable = nil
	return err
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
	cColumn := C.grn_column_create(cCtx, t.cTable, cName, C.uint(cNameLen),
		cPath, C.grn_obj_flags(flags), columnTypeObj)
	if cColumn == nil {
		return nil, errorFromRc(cCtx.rc)
	}
	column := &Column{table: t, cColumn: cColumn}
	t.addColumnToMap(name, column)
	return column, nil
}

func (t *Table) OpenColumn(name string) (*Column, error) {
	var cName *C.char
	var cNameLen C.size_t
	if name != "" {
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
		cNameLen = C.strlen(cName)
	}

	cCtx := t.db.context.cCtx
	cColumn := C.grn_obj_column(cCtx, t.cTable, cName, C.uint(cNameLen))
	if cColumn == nil {
		if cCtx.rc != SUCCESS {
			return nil, errorFromRc(cCtx.rc)
		}
		return nil, NotFoundError
	}
	column := &Column{table: t, cColumn: cColumn}
	t.addColumnToMap(name, column)
	return column, nil
}

func (t *Table) addColumnToMap(name string, column *Column) {
	if t.columns == nil {
		t.columns = make(map[string]*Column)
	}
	t.columns[name] = column
}
