package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"

func Init() error {
	rc := C.grn_init()
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func Fin() error {
	rc := C.grn_fin()
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func FinDefer(err *error) {
	err2 := Fin()
	if err2 != nil && *err == nil {
		*err = err2
	}
}
