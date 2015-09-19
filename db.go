package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

type ID int

type DB struct {
	context *Context
	cDB     *C.grn_obj
}

func (d *DB) Path() string {
	return objPath(d.context.cCtx, d.cDB)
}

func (d *DB) Close() {
	if d.cDB == nil {
		return
	}
	unlinkObj(d.context.cCtx, d.cDB)
	d.cDB = nil
}

func (d *DB) Remove() error {
	if d.cDB == nil {
		return InvalidArgumentError
	}
	err := removeObj(d.context.cCtx, d.cDB)
	d.cDB = nil
	return err
}

func (d *DB) CreateTable(name, path string, flags int, keyType *Obj) (*Table, error) {
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

	cCtx := d.context.cCtx
	if keyType == nil {
		return nil, InvalidArgumentError
	}
	cKeyType := keyType.cObj
	cTable := C.grn_table_create(cCtx, cName, C.uint(cNameLen),
		cPath, C.grn_obj_flags(flags), cKeyType, nil)
	if cTable == nil {
		return nil, errorFromRc(cCtx.rc)
	}
	return &Table{&Records{db: d, cRecords: cTable}}, nil
}

func (d *DB) OpenTable(name string) (*Table, error) {
	var cName *C.char
	var cNameLen C.size_t
	if name != "" {
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
		cNameLen = C.strlen(cName)
	}
	cCtx := d.context.cCtx
	cTable := C.grn_ctx_get(cCtx, cName, C.int(cNameLen))
	if cTable == nil {
		if cCtx.rc != SUCCESS {
			return nil, errorFromRc(cCtx.rc)
		}
		return nil, NotFoundError
	}
	return &Table{&Records{db: d, cRecords: cTable}}, nil
}

func (d *DB) OpenOrCreateTable(name, path string, flags int, keyType *Obj) (*Table, error) {
	table, err := d.OpenTable(name)
	if err != nil {
		if err != NotFoundError {
			return nil, err
		}
		table, err = d.CreateTable(name, path, flags, keyType)
	}
	return table, err
}
