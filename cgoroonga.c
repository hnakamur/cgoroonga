#include <stdio.h>
#include "cgoroonga.h"

GRN_API void cgoroonga_str_array_set(char **array, int i, char *elem) {
	array[i] = elem;
}

GRN_API char *cgoroonga_str_array_get(char **array, int i) {
	return array[i];
}

GRN_API void cgoroonga_uint_array_set(unsigned int *array, int i, unsigned int elem) {
	array[i] = elem;
}

GRN_API unsigned int cgoroonga_uint_array_get(unsigned int *array, int i) {
	return array[i];
}


GRN_API char *cgoroonga_bulk_head(grn_obj *bulk) {
	return GRN_BULK_HEAD(bulk);
}

GRN_API void cgoroonga_bulk_rewind(grn_obj *bulk) {
	GRN_BULK_REWIND(bulk);
}

GRN_API int cgoroonga_bulk_vsize(grn_obj *bulk) {
	return GRN_BULK_VSIZE(bulk);
}


GRN_API long long int cgoroonga_int64_value(grn_obj *obj) {
	return GRN_INT64_VALUE(obj);
}


GRN_API grn_snip_mapping *cgoroonga_mapping_html_escape() {
	return GRN_SNIP_MAPPING_HTML_ESCAPE;
}


GRN_API void cgoroonga_record_init(grn_obj *obj, unsigned char impl_flags, grn_id domain) {
	GRN_RECORD_INIT(obj, impl_flags, domain);
}


GRN_API void cgoroonga_text_init(grn_obj *obj, unsigned char impl_flags) {
	GRN_TEXT_INIT(obj, impl_flags);
}

GRN_API void cgoroonga_text_put(grn_ctx *ctx, grn_obj *obj, const char *str, unsigned int len) {
	GRN_TEXT_PUT(ctx, obj, str, len);
}


GRN_API void cgoroonga_time_init(grn_obj *obj, unsigned char impl_flags) {
	GRN_TIME_INIT(obj, impl_flags);
}

GRN_API void cgoroonga_time_set(grn_ctx *ctx, grn_obj *obj, long long int unix_usec) {
	GRN_TIME_SET(ctx, obj, unix_usec);
}
