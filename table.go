package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

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
