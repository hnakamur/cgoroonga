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

/*
func (t *Table) Column(name string) *Column {
	if t.columns == nil {
		return nil
	}
	return t.columns[name]
}
*/

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

func (t *Table) RecordCount() (uint, error) {
	cCtx := t.db.context.cCtx
	count := C.grn_table_size(cCtx, t.cTable)
	if cCtx.rc != SUCCESS {
		return 0, errorFromRc(cCtx.rc)
	}
	return uint(count), nil
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

func (t *Table) AddRecord(key string) (recordID ID, added bool, err error) {
	var cKey *C.char
	var cKeyLen C.size_t
	if key != "" {
		cKey = C.CString(key)
		defer C.free(unsafe.Pointer(cKey))
		cKeyLen = C.strlen(cKey)
	}
	var cAdded C.int
	recordID = ID(C.grn_table_add(
		t.db.context.cCtx,
		t.cTable,
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

func (t *Table) CreateQuery(name string) (*Expr, error) {
	expr, err := t.db.context.CreateExpr(name)
	if err != nil {
		return nil, err
	}

	cCtx := t.db.context.cCtx
	cVar, err := expr.addVar("")
	if err != nil {
		return nil, err
	}
	C.cgoroonga_record_init(cVar, 0, C.grn_obj_id(cCtx, t.cTable))
	return expr, nil
}

func (t *Table) Select(expr *Expr, res *Records, op int) (*Records, error) {
	cCtx := t.db.context.cCtx

	var cExpr *C.grn_obj
	if expr != nil {
		cExpr = expr.cExpr
	}
	var cRes *C.grn_obj
	if res != nil {
		cRes = res.cTable
	}
	cRes = C.grn_table_select(cCtx, t.cTable, cExpr, cRes, C.grn_operator(op))
	if cCtx.rc != SUCCESS {
		return nil, errorFromRc(cCtx.rc)
	}
	return &Records{&Table{db: t.db, cTable: cRes}}, nil
}
