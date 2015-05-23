package goroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

type Context struct {
	cCtx      *C.grn_ctx
	currentDB *DB
}

func NewContext() (*Context, error) {
	cCtx := C.grn_ctx_open(0)
	if cCtx == nil {
		return nil, UnknownError
	}
	return &Context{cCtx: cCtx}, nil
}

func (c *Context) Close() error {
	err := c.setCurrentDB(nil)
	if err != nil {
		return err
	}

	if c.cCtx == nil {
		return nil
	}
	rc := C.grn_ctx_close(c.cCtx)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	c.cCtx = nil
	return nil
}

func (c *Context) CreateDB(path string) (*DB, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cDB := C.grn_db_create(c.cCtx, cPath, nil)
	if cDB == nil {
		return nil, errorFromRc(c.cCtx.rc)
	}

	db := &DB{context: c, cDB: cDB}
	err := c.setCurrentDB(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *Context) OpenDB(path string) (*DB, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	cDB := C.grn_db_open(c.cCtx, cPath)
	if cDB == nil {
		return nil, errorFromRc(c.cCtx.rc)
	}

	db := &DB{context: c, cDB: cDB}
	err := c.setCurrentDB(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *Context) CurrentDB() *DB {
	return c.currentDB
}

func (c *Context) UseDB(db *DB) error {
	if c.currentDB == db {
		return nil
	}

	rc := C.grn_ctx_use(c.cCtx, db.cDB)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}

	return c.setCurrentDB(db)
}

func (c *Context) setCurrentDB(db *DB) error {
	if c.currentDB != nil {
		err := c.currentDB.Close()
		if err != nil {
			return err
		}
	}
	c.currentDB = db
	return nil
}
