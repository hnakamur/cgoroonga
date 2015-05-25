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

func (r *Records) GetKey(id ID) string {
	cCtx := r.db.context.cCtx
	cRecords := r.cRecords
	cID := C.grn_id(id)
	length := C.grn_table_get_key(cCtx, cRecords, cID, nil, 0)
	if length == 0 {
		return ""
	}

	var buf *C.char
	buf = (*C.char)(C.malloc(C.size_t(unsafe.Sizeof(*buf)) * C.size_t(length)))
	defer C.free(unsafe.Pointer(buf))
	C.grn_table_get_key(cCtx, cRecords, cID, unsafe.Pointer(buf), length)
	return C.GoString(buf)
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

func (r *Records) OpenColumn(name string) (*Column, error) {
	var cName *C.char
	var cNameLen C.size_t
	if name != "" {
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
		cNameLen = C.strlen(cName)
	}

	cCtx := r.db.context.cCtx
	cColumn := C.grn_obj_column(cCtx, r.cRecords, cName, C.uint(cNameLen))
	if cColumn == nil {
		if cCtx.rc != SUCCESS {
			return nil, errorFromRc(cCtx.rc)
		}
		return nil, NotFoundError
	}
	column := &Column{table: &Table{r}, cColumn: cColumn}
	r.addColumnToMap(name, column)
	return column, nil
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

func (r *Records) OpenTableCursor(min, max string, offset, limit, flags int) (*TableCursor, error) {
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

	cCtx := r.db.context.cCtx
	cCursor := C.grn_table_cursor_open(
		cCtx,
		r.cRecords,
		unsafe.Pointer(cMin), C.uint(cMinLen),
		unsafe.Pointer(cMax), C.uint(cMaxLen),
		C.int(offset), C.int(limit), C.int(flags))
	if cCursor == nil {
		return nil, errorFromRc(cCtx.rc)
	}
	return &TableCursor{records: r, cCursor: cCursor}, nil
}
