#ifndef CGOROONGA_H
#define CGOROONGA_H

#include <stdlib.h>
#include <string.h>
#include <groonga/groonga.h>

GRN_API char *cgoroonga_bulk_head(grn_obj *bulk);
GRN_API void cgoroonga_bulk_rewind(grn_obj *bulk);
GRN_API int cgoroonga_bulk_vsize(grn_obj *bulk);

GRN_API void cgoroonga_text_init(grn_obj *obj, unsigned char impl_flags);
GRN_API void cgoroonga_text_put(grn_ctx *ctx, grn_obj *obj, const char *str, unsigned int len);

#endif
