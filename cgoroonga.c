#include <stdio.h>
#include "cgoroonga.h"

GRN_API char *cgoroonga_bulk_head(grn_obj *bulk) {
	return GRN_BULK_HEAD(bulk);
}

GRN_API void cgoroonga_bulk_rewind(grn_obj *bulk) {
	GRN_BULK_REWIND(bulk);
}

GRN_API int cgoroonga_bulk_vsize(grn_obj *bulk) {
	return GRN_BULK_VSIZE(bulk);
}


GRN_API void cgoroonga_text_init(grn_obj *obj, unsigned char impl_flags) {
	GRN_TEXT_INIT(obj, impl_flags);
}

GRN_API void cgoroonga_text_put(grn_ctx *ctx, grn_obj *obj, const char *str, unsigned int len) {
	GRN_TEXT_PUT(ctx, obj, str, len);
}
