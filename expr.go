package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

type Expr struct {
	context *Context
	cExpr   *C.grn_obj
	cVars   []*C.grn_obj
}

func (e *Expr) Close() {
	cCtx := e.context.cCtx
	if e.cVars != nil {
		for _, cVar := range e.cVars {
			unlinkObj(cCtx, cVar)
		}
	}

	if e.cExpr != nil {
		unlinkObj(cCtx, e.cExpr)
		e.cExpr = nil
	}
}

func (e *Expr) Parse(str string, defaultColumn *Column, defaultMode, defaultOp, flags int) error {
	var cStr *C.char
	var cStrLen C.size_t
	if str != "" {
		cStr = C.CString(str)
		defer C.free(unsafe.Pointer(cStr))
		cStrLen = C.strlen(cStr)
	}

	var cDefaultColumn *C.grn_obj
	if defaultColumn != nil {
		cDefaultColumn = defaultColumn.cColumn
	}
	cCtx := e.context.cCtx
	rc := C.grn_expr_parse(cCtx, e.cExpr, cStr, C.uint(cStrLen),
		cDefaultColumn, C.grn_operator(defaultMode),
		C.grn_operator(defaultOp), C.grn_expr_flags(flags))
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (e *Expr) addVar(name string) (*C.grn_obj, error) {
	var cName *C.char
	var cNameLen C.size_t
	if name != "" {
		cName = C.CString(name)
		defer C.free(unsafe.Pointer(cName))
		cNameLen = C.strlen(cName)
	}

	cCtx := e.context.cCtx
	cVar := C.grn_expr_add_var(cCtx, e.cExpr, cName, C.uint(cNameLen))
	if cVar == nil {
		return nil, errorFromRc(cCtx.rc)
	}
	e.appendVar(cVar)
	return cVar, nil
}

func (e *Expr) appendVar(cVar *C.grn_obj) {
	if e.cVars == nil {
		e.cVars = make([]*C.grn_obj, 1)
	}
	e.cVars = append(e.cVars, cVar)
}

func (e *Expr) AppendOp(op, nargs int) error {
	cCtx := e.context.cCtx
	rc := C.grn_expr_append_op(cCtx, e.cExpr, C.grn_operator(op), C.int(nargs))
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}
