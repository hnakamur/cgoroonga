package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import (
	"time"
	"unsafe"
)

func TimeInit(obj *Obj, implFlags int) {
	C.go_grn_time_init(
		(*C.struct__grn_obj)(unsafe.Pointer(obj)),
		C.uchar(implFlags))
}

func (c *Ctx) TimeSet(obj *Obj, t time.Time) {
	C.go_grn_time_set(
		(*C.struct__grn_ctx)(unsafe.Pointer(c)),
		(*C.struct__grn_obj)(unsafe.Pointer(obj)),
		C.longlong(t.UnixNano()/1000))
}

func TimeValue(obj *Obj) time.Time {
	usec := C.go_grn_time_value(
		(*C.struct__grn_obj)(unsafe.Pointer(obj)))
	return time.Unix(int64(usec/1000000), int64((usec%1000000)*1000))
}
