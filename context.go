package cgoroonga

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
	c.setCurrentDB(nil)

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
	c.setCurrentDB(db)
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
	c.setCurrentDB(db)
	return db, nil
}

func (c *Context) OpenOrCreateDB(path string) (*DB, error) {
	db, err := c.OpenDB(path)
	if err != nil {
		if err != NoSuchFileOrDirectoryError {
			return nil, err
		}
		db, err = c.CreateDB(path)
	}
	return db, err
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

	c.setCurrentDB(db)
	return nil
}

func (c *Context) setCurrentDB(db *DB) {
	if c.currentDB != nil {
		c.currentDB.Close()
	}
	c.currentDB = db
}

func (c *Context) CreateExpr(name string) (*Expr, error) {
	var cName *C.char
	var cNameLen C.size_t
	if name != "" {
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
		cNameLen = C.strlen(cName)
	}

	cExpr := C.grn_expr_create(c.cCtx, cName, C.uint(cNameLen))
	if cExpr == nil {
		return nil, errorFromRc(c.cCtx.rc)
	}
	return &Expr{context: c, cExpr: cExpr}, nil
}

func grnCtxGet(cCtx *C.grn_ctx, name string) *C.grn_obj {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	cNameLen := C.strlen(cName)

	return C.grn_ctx_get(cCtx, cName, C.int(cNameLen))
}
