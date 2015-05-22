package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func (c *Ctx) TableOpen(name string) (*Obj, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cNameLen := C.strlen(cName)
	table := C.grn_ctx_get(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)), cName, C.int(cNameLen))
	if table == nil {
		return nil, NoSuchFileOrDirectoryError
	}
	return (*Obj)(unsafe.Pointer(table)), nil
}

func (c *Ctx) TableOpenOrCreate(name string, path string, flags ObjFlags, keyType, valueType *Obj) (*Obj, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cNameLen := C.strlen(cName)
	table := C.grn_ctx_get(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)), cName, C.int(cNameLen))
	if table == nil {
		var cPath *C.char
		if path != "" {
			cPath = C.CString(path)
			defer C.free(unsafe.Pointer(cPath))
		}

		table = C.grn_table_create(
			(*C.struct__grn_ctx)(unsafe.Pointer(c)),
			cName, C.uint(cNameLen),
			cPath,
			C.grn_obj_flags(flags),
			(*C.struct__grn_obj)(unsafe.Pointer(keyType)),
			(*C.struct__grn_obj)(unsafe.Pointer(valueType)))
		if table == nil {
			return nil, TableCreateError
		}
	}
	return (*Obj)(unsafe.Pointer(table)), nil
}

func (c *Ctx) TableAdd(table *Obj, key string) (recordID ID, added bool, err error) {
	var cKey *C.char
	var cKeyLen C.size_t
	if key != "" {
		cKey = C.CString(key)
		defer C.free(unsafe.Pointer(cKey))

		cKeyLen = C.strlen(cKey)
	}
	var cAdded C.int
	recordID = ID(C.grn_table_add(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(table)),
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

func (c *Ctx) TableSelect(table, expr, res *Obj, op Operator) (*Obj, error) {
	result := C.grn_table_select(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(table)),
		(*C.struct__grn_obj)(unsafe.Pointer(expr)),
		(*C.struct__grn_obj)(unsafe.Pointer(res)),
		C.grn_operator(op))
	if result == nil {
		return nil, errorFromRc(c.rc)
	}
	return (*Obj)(unsafe.Pointer(result)), nil
}

func (c *Ctx) TableSize(table *Obj) (uint, error) {
	n := C.grn_table_size(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(table)))
	if c.rc != SUCCESS {
		return 0, errorFromRc(c.rc)
	}
	return uint(n), nil
}
