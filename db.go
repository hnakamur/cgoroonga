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
		return nil, DBCreateError
	}
	return db, nil
}

func (c *Ctx) WithDB(path string, optarg *CreateOptArg, handler func(ctx *Ctx, db *Obj) error) (err error) {
	var db *Obj
	db, err = c.DBOpenOrCreate(path, optarg)
	if err != nil {
		return
	}
	defer func() {
		err2 := c.ObjClose(db)
		if err2 != nil && err == nil {
			err = err2
		}
	}()

	return handler(c, db)
}
