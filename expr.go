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

func (e *Expr) Parse(str string, defaultColumn *Expr, defaultMode, defaultOp, flags int) error {
	var cStr *C.char
	var cStrLen C.size_t
	if str != "" {
		cStr = C.CString(str)
		defer C.free(unsafe.Pointer(cStr))
		cStrLen = C.strlen(cStr)
	}

	var cDefaultColumn *C.grn_obj
	if defaultColumn != nil {
		cDefaultColumn = defaultColumn.cExpr
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

func (e *Expr) Snippet(flags, width, maxResults int, htmlEscape bool, tagPairs [][]string) *Snippet {
	var cTag *C.char
	var cTagLen C.uint
	n := len(tagPairs)
	openTags := (**C.char)(C.malloc(C.size_t(unsafe.Sizeof(cTag)) * C.size_t(n)))
	defer C.free(unsafe.Pointer(openTags))
	openTagLens := (*C.uint)(C.malloc(C.size_t(unsafe.Sizeof(cTagLen)) * C.size_t(n)))
	defer C.free(unsafe.Pointer(openTagLens))
	closeTags := (**C.char)(C.malloc(C.size_t(unsafe.Sizeof(cTag)) * C.size_t(n)))
	defer C.free(unsafe.Pointer(closeTags))
	closeTagLens := (*C.uint)(C.malloc(C.size_t(unsafe.Sizeof(cTagLen)) * C.size_t(n)))
	defer C.free(unsafe.Pointer(closeTagLens))
	for i := 0; i < n; i++ {
		tagPair := tagPairs[i]
		openTag := C.CString(tagPair[0])
		closeTag := C.CString(tagPair[1])
		C.cgoroonga_str_array_set(openTags, C.int(i), openTag)
		C.cgoroonga_uint_array_set(openTagLens, C.int(i), C.uint(C.strlen(openTag)))
		C.cgoroonga_str_array_set(closeTags, C.int(i), closeTag)
		C.cgoroonga_uint_array_set(closeTagLens, C.int(i), C.uint(C.strlen(closeTag)))
	}
	defer func() {
		for i := 0; i < n; i++ {
			C.free(unsafe.Pointer(C.cgoroonga_str_array_get(openTags, C.int(i))))
			C.free(unsafe.Pointer(C.cgoroonga_str_array_get(closeTags, C.int(i))))
		}
	}()
	var mapping *C.grn_snip_mapping
	if htmlEscape {
		mapping = C.cgoroonga_mapping_html_escape()
	}
	cSnip := C.grn_expr_snip(e.context.cCtx, e.cExpr, C.int(flags),
		C.uint(width), C.uint(maxResults), C.uint(n),
		openTags, openTagLens, closeTags, closeTagLens, mapping)
	return &Snippet{context: e.context, cSnip: cSnip}
}
