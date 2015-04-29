#include "cgoroonga.h"

GRN_API grn_obj *go_grn_db_open_or_create(grn_ctx *ctx, const char *path, grn_db_create_optarg *optarg) {
	grn_obj *db;
	db = grn_db_open(ctx, path);
	if (!db) {
		db = grn_db_create(ctx, path, optarg);
	}
	return db;
}

GRN_API void go_grn_text_init(grn_obj *text, unsigned char impl_flags) {
	GRN_TEXT_INIT(text, impl_flags);
}

GRN_API grn_rc go_grn_text_put(grn_ctx *ctx, grn_obj *bulk, const char *str, unsigned int len) {
	return grn_bulk_write(ctx, bulk, str, len);
}

GRN_API void go_grn_bulk_rewind(grn_obj *bulk) {
	GRN_BULK_REWIND(bulk);
}

GRN_API char *go_grn_bulk_head(grn_obj *bulk) {
	return GRN_BULK_HEAD(bulk);
}

GRN_API void go_grn_record_init(grn_obj *obj, unsigned char flags, grn_id domain) {
	GRN_VALUE_FIX_SIZE_INIT(obj, flags, domain);
}
