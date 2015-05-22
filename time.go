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

func (c *Ctx) GetTime(column *Obj, recordID ID) time.Time {
	var buf Obj
	TimeInit(&buf, 0)
	defer c.ObjUnlink(&buf)
	c.ObjGetValue(column, recordID, &buf)
	return TimeValue(&buf)
}

func (c *Ctx) SetTime(column *Obj, recordID ID, t time.Time) error {
	var value Obj
	TimeInit(&value, 0)
	defer c.ObjUnlink(&value)
	c.TimeSet(&value, t)
	return c.ObjSetValue(column, recordID, &value, OBJ_SET)
}
