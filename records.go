package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

type Records struct {
	db       *DB
	cRecords *C.grn_obj
	columns  map[string]*Column
}

func (r *Records) Column(name string) *Column {
	if r.columns == nil {
		return nil
	}
	return r.columns[name]
}

func (r *Records) RecordCount() (uint, error) {
	cCtx := r.db.context.cCtx
	count := C.grn_table_size(cCtx, r.cRecords)
	if cCtx.rc != SUCCESS {
		return 0, errorFromRc(cCtx.rc)
	}
	return uint(count), nil
}

func (r *Records) Path() string {
	return objPath(r.db.context.cCtx, r.cRecords)
}

func (r *Records) Name() string {
	return objName(r.db.context.cCtx, r.cRecords)
}

func (r *Records) Close() {
	for name, column := range r.columns {
		column.Close()
		delete(r.columns, name)
	}

	if r.cRecords == nil {
		return
	}
	unlinkObj(r.db.context.cCtx, r.cRecords)
	r.cRecords = nil
}

func (r *Records) Remove() error {
	if r.cRecords == nil {
		return InvalidArgumentError
	}
	err := removeObj(r.db.context.cCtx, r.cRecords)
	r.cRecords = nil
	return err
}

func (r *Records) addColumnToMap(name string, column *Column) {
	if r.columns == nil {
		r.columns = make(map[string]*Column)
	}
	r.columns[name] = column
}

func (r *Records) AddRecord(key string) (recordID ID, added bool, err error) {
	var cKey *C.char
	var cKeyLen C.size_t
	if key != "" {
		cKey = C.CString(key)
		defer C.free(unsafe.Pointer(cKey))
		cKeyLen = C.strlen(cKey)
	}
	var cAdded C.int
	recordID = ID(C.grn_table_add(
		r.db.context.cCtx,
		r.cRecords,
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

func (r *Records) CreateQuery(name string) (*Expr, error) {
	expr, err := r.db.context.CreateExpr(name)
	if err != nil {
		return nil, err
	}

	cCtx := r.db.context.cCtx
	cVar, err := expr.addVar("")
	if err != nil {
		return nil, err
	}
	C.cgoroonga_record_init(cVar, 0, C.grn_obj_id(cCtx, r.cRecords))
	return expr, nil
}

func (r *Records) Select(expr *Expr, records *Records, op int) (*Records, error) {
	cCtx := r.db.context.cCtx

	var cExpr *C.grn_obj
	if expr != nil {
		cExpr = expr.cExpr
	}
	var cRecords *C.grn_obj
	if records != nil {
		cRecords = records.cRecords
	}
	cRecords = C.grn_table_select(cCtx, r.cRecords, cExpr, cRecords,
		C.grn_operator(op))
	if cCtx.rc != SUCCESS {
		return nil, errorFromRc(cCtx.rc)
	}
	return &Records{db: r.db, cRecords: cRecords}, nil
}
