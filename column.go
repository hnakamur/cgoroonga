package cgoroonga

/*
#cgo LDFLAGS: -lgroonga
#include "cgoroonga.h"
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

type Column struct {
	table   *Table
	cColumn *C.grn_obj
}

func (c *Column) Name() string {
	cCtx := c.table.db.context.cCtx
	length := C.grn_column_name(cCtx, c.cColumn, nil, 0)
	if length == 0 {
		return ""
	}

	var buf *C.char
	buf = (*C.char)(C.malloc(C.size_t(unsafe.Sizeof(*buf)) * C.size_t(length)))
	defer C.free(unsafe.Pointer(buf))
	C.grn_column_name(cCtx, c.cColumn, buf, length)
	return C.GoString(buf)
}

func (c *Column) Path() string {
	return objPath(c.table.db.context.cCtx, c.cColumn)
}

func (c *Column) DataType() ID {
	cCtx := c.table.db.context.cCtx
	return ID(C.grn_obj_get_range(cCtx, c.cColumn))
}

func (c *Column) Close() {
	if c.cColumn == nil {
		return
	}
	unlinkObj(c.table.db.context.cCtx, c.cColumn)
	c.cColumn = nil
}

func (c *Column) Remove() error {
	if c.cColumn == nil {
		return InvalidArgumentError
	}
	err := removeObj(c.table.db.context.cCtx, c.cColumn)
	c.cColumn = nil
	return err
}

func (c *Column) SetString(id ID, str string) error {
	var cStr *C.char
	var cStrLen C.size_t
	if str != "" {
		cStr = C.CString(str)
		defer C.free(unsafe.Pointer(cStr))
		cStrLen = C.strlen(cStr)
	}

	cCtx := c.table.db.context.cCtx
	var buf C.grn_obj
	c.initTextObj(&buf)
	defer unlinkObj(cCtx, &buf)
	C.cgoroonga_text_put(cCtx, &buf, cStr, C.uint(cStrLen))
	rc := C.grn_obj_set_value(cCtx, c.cColumn, C.grn_id(id), &buf,
		C.GRN_OBJ_SET)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (c *Column) GetString(id ID) string {
	cCtx := c.table.db.context.cCtx
	var buf C.grn_obj
	c.initTextObj(&buf)
	defer unlinkObj(cCtx, &buf)
	C.grn_obj_get_value(cCtx, c.cColumn, C.grn_id(id), &buf)
	return C.GoStringN(C.cgoroonga_bulk_head(&buf),
		C.cgoroonga_bulk_vsize(&buf))
}

func (c *Column) initTextObj(buf *C.grn_obj) {
	cCtx := c.table.db.context.cCtx
	range_ := C.grn_obj_get_range(cCtx, c.cColumn)
	type_ := C.cgoroonga_obj_header_type(c.cColumn)
	switch {
	case type_ == C.GRN_COLUMN_VAR_SIZE || type_ == C.GRN_COLUMN_FIX_SIZE:
		switch C.cgoroonga_obj_header_flags(c.cColumn) & C.GRN_OBJ_COLUMN_TYPE_MASK {
		case C.GRN_OBJ_COLUMN_VECTOR:
			cRangeObj := C.grn_ctx_at(cCtx, range_)
			if C.cgoroonga_obj_tablep(c.cColumn) != 0 ||
				(C.cgoroonga_obj_header_flags(cRangeObj)&C.GRN_OBJ_KEY_VAR_SIZE) == 0 {
				C.cgoroonga_obj_init(buf, C.GRN_UVECTOR, 0, range_)
			} else {
				C.cgoroonga_obj_init(buf, C.GRN_VECTOR, 0, range_)
			}
		case C.GRN_OBJ_COLUMN_SCALAR:
			C.cgoroonga_obj_init(buf, C.GRN_BULK, 0, range_)
		default:
			panic(fmt.Sprintf("unsupported column flags: %x", C.cgoroonga_obj_header_flags(c.cColumn)))
		}
	case type_ == C.GRN_COLUMN_INDEX:
		// Do nothing
	case type_ == C.GRN_ACCESSOR:
		C.cgoroonga_obj_init(buf, C.GRN_BULK, 0, range_)
	default:
		panic(fmt.Sprintf("unsupported header type: %x", type_))
	}

	// switch (column->header.type) {
	// case GRN_COLUMN_VAR_SIZE:
	// case GRN_COLUMN_FIX_SIZE:
	//   switch (column->header.flags & GRN_OBJ_COLUMN_TYPE_MASK) {
	//   case GRN_OBJ_COLUMN_VECTOR:
	//     dump_record_column_vector(ctx, outbuf, id, column, range, &buf);
	//     break;
	//   case GRN_OBJ_COLUMN_SCALAR:
	//     {
	//       GRN_OBJ_INIT(&buf, GRN_BULK, 0, range);
	//       grn_obj_get_value(ctx, column, id, &buf);
	//       grn_text_otoj(ctx, outbuf, &buf, NULL);
	//       grn_obj_unlink(ctx, &buf);
	//     }
	//     break;
	//   default:
	//     ERR(GRN_OPERATION_NOT_SUPPORTED,
	//         "unsupported column type: %#x",
	//         column->header.type);
	//     break;
	//   }
	//   break;
	// case GRN_COLUMN_INDEX:
	//   break;
	// case GRN_ACCESSOR:
	//   {
	//     GRN_OBJ_INIT(&buf, GRN_BULK, 0, range);
	//     grn_obj_get_value(ctx, column, id, &buf);
	//     /* XXX maybe, grn_obj_get_range() should not unconditionally return
	//        GRN_DB_INT32 when column is GRN_ACCESSOR and
	//        GRN_ACCESSOR_GET_VALUE */
	//     if (is_value_column) {
	//       buf.header.domain = ((grn_db_obj *)table)->range;
	//     }
	//     grn_text_otoj(ctx, outbuf, &buf, NULL);
	//     grn_obj_unlink(ctx, &buf);
	//   }
	//   break;
	// default:
	//   ERR(GRN_OPERATION_NOT_SUPPORTED,
	//       "unsupported header type %#x",
	//       column->header.type);
	//   break;
	// }
}

func (c *Column) SetTime(id ID, t time.Time) error {
	cCtx := c.table.db.context.cCtx
	var buf C.grn_obj
	C.cgoroonga_time_init(&buf, 0)
	defer unlinkObj(cCtx, &buf)
	// convert nano seconds to micro seconds
	usec := t.UnixNano() / 1000
	C.cgoroonga_time_set(cCtx, &buf, C.longlong(usec))
	rc := C.grn_obj_set_value(cCtx, c.cColumn, C.grn_id(id), &buf,
		C.GRN_OBJ_SET)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (c *Column) GetTime(id ID) time.Time {
	cCtx := c.table.db.context.cCtx
	var buf C.grn_obj
	C.cgoroonga_time_init(&buf, 0)
	defer unlinkObj(cCtx, &buf)
	C.grn_obj_get_value(cCtx, c.cColumn, C.grn_id(id), &buf)
	vsize := C.cgoroonga_bulk_vsize(&buf)
	if vsize == 0 {
		return time.Unix(0, 0)
	}
	usec := C.cgoroonga_int64_value(&buf)
	return time.Unix(int64(usec/1000000), int64((usec%1000000)*1000))
}

func (c *Column) SetStringArray(id ID, values []string) error {
	cCtx := c.table.db.context.cCtx
	var vec C.grn_obj
	c.initTextObj(&vec)
	defer unlinkObj(cCtx, &vec)

	typeID := C.grn_obj_get_range(cCtx, c.cColumn)
	for _, value := range values {
		var cValue *C.char
		var cValueLen C.size_t
		if value != "" {
			cValue = C.CString(value)
			cValueLen = C.strlen(cValue)
		}
		rc := C.grn_vector_add_element(cCtx, &vec, cValue, C.uint(cValueLen), 0, typeID)
		if value != "" {
			defer C.free(unsafe.Pointer(cValue))
		}
		if rc != SUCCESS {
			return errorFromRc(rc)
		}
	}

	rc := C.grn_obj_set_value(cCtx, c.cColumn, C.grn_id(id), &vec,
		C.GRN_OBJ_SET)
	if rc != SUCCESS {
		return errorFromRc(rc)
	}
	return nil
}

func (c *Column) GetStringArray(id ID) ([]string, error) {
	cCtx := c.table.db.context.cCtx
	var vec C.grn_obj
	c.initTextObj(&vec)
	defer unlinkObj(cCtx, &vec)
	C.grn_obj_get_value(cCtx, c.cColumn, C.grn_id(id), &vec)
	if cCtx.rc != SUCCESS {
		return nil, errorFromRc(cCtx.rc)
	}

	type_ := C.cgoroonga_obj_header_type(&vec)
	var strings []string
	switch type_ {
	case C.GRN_UVECTOR:
		size := C.grn_uvector_size(cCtx, &vec)
		strings = make([]string, size)
		range_ := C.grn_ctx_at(cCtx, C.cgoroonga_obj_header_domain(&vec))
		rangeHeaderType := C.cgoroonga_obj_header_type(range_)
		rangeIsType := (rangeHeaderType == C.GRN_TYPE)
		if rangeIsType {
		} else {
			var i C.uint
			for i = 0; i < size; i++ {
				sourceID := C.grn_uvector_get_element(cCtx, &vec, i, nil)
				if cCtx.rc != SUCCESS {
					return nil, errorFromRc(cCtx.rc)
				}
				//strings[i] = C.GoString(value)
			}
		}

		//if (range_is_type) {
		//  unsigned int i, n;
		//  char *raw_elements;
		//  unsigned int element_size;
		//  grn_obj element;

		//  raw_elements = GRN_BULK_HEAD(uvector);
		//  element_size = GRN_TYPE_SIZE(DB_OBJ(range));
		//  n = GRN_BULK_VSIZE(uvector) / element_size;

		//  grn_output_array_open(ctx, outbuf, output_type, "VECTOR", n);
		//  GRN_OBJ_INIT(&element, GRN_BULK, 0, uvector->header.domain);
		//  for (i = 0; i < n; i++) {
		//    GRN_BULK_REWIND(&element);
		//    grn_bulk_write_from(ctx, &element, raw_elements + (element_size * i),
		//                        0, element_size);
		//    grn_output_obj(ctx, outbuf, output_type, &element, NULL);
		//  }
		//  GRN_OBJ_FIN(ctx, &element);
		//  grn_output_array_close(ctx, outbuf, output_type);
		//} else {
		//  unsigned int i, n;
		//  grn_obj id_value;
		//  grn_obj key_value;

		//  GRN_UINT32_INIT(&id_value, 0);
		//  GRN_OBJ_INIT(&key_value, GRN_BULK, 0, range->header.domain);

		//  n = grn_vector_size(ctx, uvector);
		//  if (with_weight) {
		//    grn_output_map_open(ctx, outbuf, output_type, "WEIGHT_VECTOR", n);
		//  } else {
		//    grn_output_array_open(ctx, outbuf, output_type, "VECTOR", n);
		//  }

		//  for (i = 0; i < n; i++) {
		//    grn_id id;
		//    unsigned int weight;

		//    id = grn_uvector_get_element(ctx, uvector, i, &weight);
		//    if (range->header.type == GRN_TABLE_NO_KEY) {
		//      GRN_UINT32_SET(ctx, &id_value, id);
		//      grn_output_obj(ctx, outbuf, output_type, &id_value, NULL);
		//    } else {
		//      GRN_BULK_REWIND(&key_value);
		//      grn_table_get_key2(ctx, range, id, &key_value);
		//      grn_output_obj(ctx, outbuf, output_type, &key_value, NULL);
		//    }

		//    if (with_weight) {
		//      grn_output_uint64(ctx, outbuf, output_type, weight);
		//    }
		//  }

		//  if (with_weight) {
		//    grn_output_map_close(ctx, outbuf, output_type);
		//  } else {
		//    grn_output_array_close(ctx, outbuf, output_type);
		//  }

		//  GRN_OBJ_FIN(ctx, &id_value);
		//  GRN_OBJ_FIN(ctx, &key_value);
		//}

	case C.GRN_VECTOR:
		size := C.grn_vector_size(cCtx, &vec)
		strings = make([]string, size)
		var i C.uint
		for i = 0; i < size; i++ {
			var value *C.char
			var domain C.grn_id
			C.grn_vector_get_element(cCtx, &vec, i, &value, nil, &domain)
			if cCtx.rc != SUCCESS {
				return nil, errorFromRc(cCtx.rc)
			}
			strings[i] = C.GoString(value)
		}
	default:
		return nil, fmt.Errorf("Unsupported column type: 0x%d", type_)
	}
	return strings, nil
}
