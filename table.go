package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"

static void init_source_ids(grn_obj *obj) {
	GRN_UINT32_INIT(obj, GRN_OBJ_VECTOR);
}

static void append_source_id(grn_ctx *ctx, grn_obj *source_ids, grn_id source_id) {
	GRN_UINT32_PUT(ctx, source_ids, source_id);
}
*/
import "C"
import "unsafe"

type Table struct {
	*Records
}

func (t *Table) CreateColumn(name, path string, flags int, type_ *Obj, source ...string) (*Column, error) {
	cCtx := t.db.context.cCtx
	if type_ == nil {
		return nil, InvalidArgumentError
	}
	cType := type_.cObj
	cColumn := grnColumnCreate(cCtx, t.cRecords, name, path, flags, cType)
	if cColumn == nil {
		return nil, errorFromRc(cCtx.rc)
	}
	if len(source) > 0 {
		var sourceIDs C.grn_obj
		C.init_source_ids(&sourceIDs)
		defer C.grn_obj_unlink(cCtx, &sourceIDs)
		for _, s := range source {
			cSourceColumn := grnObjColumn(cCtx, cType, s)
			if cSourceColumn == nil {
				return nil, InvalidArgumentError
			}

			var sourceID C.grn_id
			if C.cgoroonga_obj_header_type(cSourceColumn) == C.GRN_ACCESSOR {
				if s != "_key" {
					return nil, InvalidArgumentError
				}
				sourceID = C.grn_obj_id(cCtx, t.cRecords)
			} else {
				sourceID = C.grn_obj_id(cCtx, cSourceColumn)
			}
			C.append_source_id(cCtx, &sourceIDs, sourceID)
			C.grn_obj_unlink(cCtx, cSourceColumn)
		}
		rc := C.grn_obj_set_info(cCtx, cColumn, C.GRN_INFO_SOURCE, &sourceIDs)
		if rc != SUCCESS {
			return nil, errorFromRc(rc)
		}
	}
	column := &Column{table: t, cColumn: cColumn}
	t.addColumnToMap(name, column)
	return column, nil
}

func grnColumnCreate(cCtx *C.grn_ctx, cTable *C.grn_obj, name, path string, flags int, type_ *C.grn_obj) *C.grn_obj {
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

	return C.grn_column_create(cCtx, cTable, cName, C.uint(cNameLen),
		cPath, C.grn_obj_flags(flags), type_)
}

func (t *Table) OpenOrCreateColumn(name, path string, flags int, type_ *Obj, source ...string) (*Column, error) {
	column, err := t.OpenColumn(name)
	if err != nil {
		if err != NotFoundError {
			return nil, err
		}
		column, err = t.CreateColumn(name, path, flags, type_, source...)
	}
	return column, err
}

func (t *Table) SetDefaultTokenizer(name string) error {
	cCtx := t.db.context.cCtx
	cTokenizer := grnCtxGet(cCtx, name)
	if cTokenizer == nil {
		return InvalidArgumentError
	}
	rc := C.grn_obj_set_info(cCtx, t.cRecords, C.GRN_INFO_DEFAULT_TOKENIZER, cTokenizer)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
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
