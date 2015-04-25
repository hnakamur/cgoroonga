#ifndef CGOROONGA_H
#define CGOROONGA_H

#include <stdlib.h>
#include <string.h>
#include <groonga/groonga.h>

GRN_API void go_grn_bulk_rewind(grn_obj *bulk);
GRN_API char *go_grn_bulk_head(grn_obj *bulk);

GRN_API grn_obj *go_grn_db_open_or_create(grn_ctx *ctx, const char *path, grn_db_create_optarg *optarg);

GRN_API void go_grn_text_init(grn_obj *text, unsigned char impl_flags);
GRN_API grn_rc go_grn_text_put(grn_ctx *ctx, grn_obj *bulk, const char *str, unsigned int len);

#endif
