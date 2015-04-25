package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import "unsafe"

func TextInit(text *Obj, implFlags int) {
	C.go_grn_text_init(
		(*C.struct__grn_obj)(unsafe.Pointer(text)),
		C.uchar(implFlags))
}

func (c *Ctx) TextPut(bulk *Obj, str string) error {
	var cStr *C.char
	var cStrLen C.size_t
	if str != "" {
		cStr = C.CString(str)
		defer C.free(unsafe.Pointer(cStr))

		cStrLen = C.strlen(cStr)
	}
	rc := C.go_grn_text_put(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(bulk)),
		cStr,
		C.uint(cStrLen))
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}
