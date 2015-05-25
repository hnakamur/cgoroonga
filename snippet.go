package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

type Snippet struct {
	context *Context
	cSnip   *C.grn_obj
}

func (s *Snippet) Exec(str string) ([]string, error) {
	cStr := C.CString(str)
	defer C.free(unsafe.Pointer(cStr))
	cStrLen := C.strlen(cStr)
	var nResults uint
	var maxTaggedLength uint

	cCtx := s.context.cCtx
	rc := C.grn_snip_exec(
		cCtx,
		s.cSnip,
		cStr,
		C.uint(cStrLen),
		(*C.uint)(unsafe.Pointer(&nResults)),
		(*C.uint)(unsafe.Pointer(&maxTaggedLength)))
	if rc != SUCCESS {
		return nil, errorFromRc(rc)
	}
	results := []string{}
	var c C.char
	buf := (*C.char)(C.malloc(C.size_t(unsafe.Sizeof(c)) * C.size_t(maxTaggedLength)))
	defer C.free(unsafe.Pointer(buf))
	var i uint
	for i = 0; i < nResults; i++ {
		var resultLen C.uint
		rc := C.grn_snip_get_result(
			cCtx,
			s.cSnip,
			C.uint(i),
			buf,
			&resultLen)
		if rc != SUCCESS {
			return nil, errorFromRc(rc)
		}
		results = append(results, C.GoString(buf))
	}
	return results, nil
}
