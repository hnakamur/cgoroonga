package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func (c *Ctx) ColumnOpen(table *Obj, name string) (*Obj, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cNameLen := C.strlen(cName)
	column := C.grn_obj_column(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(table)),
		cName, C.uint(cNameLen))
	if column == nil {
		return nil, NoSuchFileOrDirectoryError
	}
	return (*Obj)(unsafe.Pointer(column)), nil
}

func (c *Ctx) ColumnOpenOrCreate(table *Obj, name string, path string, flags ObjFlags, columnType *Obj) (*Obj, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	cNameLen := C.strlen(cName)
	column := C.grn_obj_column(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(table)),
		cName, C.uint(cNameLen))
	if column == nil {
		var cPath *C.char
		if path != "" {
			cPath = C.CString(path)
			defer C.free(unsafe.Pointer(cPath))
		}

		column = C.grn_column_create(
			(*C.struct__grn_ctx)(unsafe.Pointer(c)),
			(*C.struct__grn_obj)(unsafe.Pointer(table)),
			cName, C.uint(cNameLen),
			cPath,
			C.grn_obj_flags(flags),
			(*C.struct__grn_obj)(unsafe.Pointer(columnType)))
		if column == nil {
			return nil, ColumnCreateError
		}
	}
	return (*Obj)(unsafe.Pointer(column)), nil
}
