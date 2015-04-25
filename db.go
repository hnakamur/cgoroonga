package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func (c *Ctx) DBOpenOrCreate(path string, optarg *CreateOptArg) (*Obj, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	db := (*Obj)(unsafe.Pointer(
		C.go_grn_db_open_or_create(
			(*C.struct__grn_ctx)(unsafe.Pointer(c)),
			cPath,
			(*C.struct__grn_db_create_optarg)(unsafe.Pointer(optarg)),
		),
	))
	if db == nil {
		return nil, Error(UNKNOWN_ERROR) //TODO: change error code
	}
	return db, nil
}
