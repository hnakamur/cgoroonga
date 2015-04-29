package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func (c *Ctx) ExprParse(expr *Obj, str string, defaultColumn *Obj, defaultMode, defaultOp Operator, flags ExprFlags) error {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))

	cStrLen := C.strlen(cStr)
	rc := C.grn_expr_parse(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(expr)),
		cStr, C.uint(cStrLen),
		(*C.struct__grn_obj)(unsafe.Pointer(defaultColumn)),
		C.grn_operator(defaultMode),
		C.grn_operator(defaultOp),
		C.grn_expr_flags(flags))
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}
