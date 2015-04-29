package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"

type Ctx C.struct__grn_ctx
type Obj C.struct__grn_obj
type TableCursor Obj
type CreateOptArg C.struct__grn_db_create_optarg
type ObjFlags int
type ID int
type Operator int
