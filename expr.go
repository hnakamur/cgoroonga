package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func (c *Ctx) ExprParse(expr *Obj, str string, defaultColumn *Obj, defaultMode, defaultOp Operator, flags int) error {
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

func (c *Ctx) ExprAppendOp(expr *Obj, op Operator, nargs int) error {
	rc := C.grn_expr_append_op(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(expr)),
		C.grn_operator(op),
		C.int(nargs))
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (c *Ctx) ExprCreateForQuery(table *Obj) (expr, var_ *Obj, err error) {
	cCtx := (*C.struct__grn_ctx)(unsafe.Pointer(c))
	cExpr := C.grn_expr_create(cCtx, nil, 0)
	var cVar *C.grn_obj
	if cExpr != nil {
		cVar = C.grn_expr_add_var(cCtx, cExpr, nil, 0)
		if cVar != nil {
			cTable := (*C.struct__grn_obj)(unsafe.Pointer(table))
			C.go_grn_record_init(cVar, 0, C.grn_obj_id(cCtx, cTable))
		}
	}
	if cExpr == nil || cVar == nil {
		return nil, nil, errorFromRc(c.rc)
	}
	return (*Obj)(cExpr), (*Obj)(cVar), nil
}

func (c *Ctx) ExprSnippet(expr *Obj, flags, width, maxResults int, htmlEscape bool, tagPairs [][]string) *Obj {
	n := len(tagPairs)
	openTags := C.go_grn_alloc_str_array(C.int(n))
	defer C.free(unsafe.Pointer(openTags))
	openTagLens := C.go_grn_alloc_uint_array(C.int(n))
	defer C.free(unsafe.Pointer(openTagLens))
	closeTags := C.go_grn_alloc_str_array(C.int(n))
	defer C.free(unsafe.Pointer(closeTags))
	closeTagLens := C.go_grn_alloc_uint_array(C.int(n))
	defer C.free(unsafe.Pointer(closeTagLens))
	for i := 0; i < n; i++ {
		tagPair := tagPairs[i]
		openTag := C.CString(tagPair[0])
		closeTag := C.CString(tagPair[1])
		C.go_grn_str_array_set(openTags, C.int(i), openTag)
		C.go_grn_uint_array_set(openTagLens, C.int(i),
			C.uint(C.strlen(openTag)))
		C.go_grn_str_array_set(closeTags, C.int(i), closeTag)
		C.go_grn_uint_array_set(closeTagLens, C.int(i),
			C.uint(C.strlen(closeTag)))
	}
	defer C.go_grn_str_array_free_elems(openTags, C.int(n))
	defer C.go_grn_str_array_free_elems(closeTags, C.int(n))
	var mapping *C.grn_snip_mapping
	if htmlEscape {
		mapping = C.go_grn_mapping_html_escape()
	}
	return (*Obj)(C.grn_expr_snip(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(expr)),
		C.int(flags), C.uint(width), C.uint(maxResults),
		C.uint(n),
		openTags, openTagLens,
		closeTags, closeTagLens,
		mapping))
}

func (c *Ctx) SnipExec(snip *Obj, str string) ([]string, error) {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))
	cStrLen := C.strlen(cStr)
	var nResults uint
	var maxTaggedLength uint

	rc := C.grn_snip_exec(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(snip)),
		cStr,
		C.uint(cStrLen),
		(*C.uint)(unsafe.Pointer(&nResults)),
		(*C.uint)(unsafe.Pointer(&maxTaggedLength)))
	if rc != SUCCESS {
		return nil, errorFromRc(rc)
	}
	results := []string{}
	var i uint
	for i = 0; i < nResults; i++ {
		buf := C.go_grn_malloc_str(C.int(maxTaggedLength))
		var resultLen C.uint
		rc := C.grn_snip_get_result(
			(*C.struct__grn_ctx)(unsafe.Pointer(c)),
			(*C.struct__grn_obj)(unsafe.Pointer(snip)),
			C.uint(i),
			buf,
			&resultLen)
		if rc != SUCCESS {
			C.free(unsafe.Pointer(buf))
			return nil, errorFromRc(rc)
		}
		results = append(results, C.GoString(buf))
		C.free(unsafe.Pointer(buf))
	}

	return results, nil
}
