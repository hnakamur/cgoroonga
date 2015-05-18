#include <stdio.h>
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

GRN_API int go_grn_bulk_vsize(grn_obj *bulk) {
	return GRN_BULK_VSIZE(bulk);
}

GRN_API void go_grn_record_init(grn_obj *obj, unsigned char flags, grn_id domain) {
	GRN_VALUE_FIX_SIZE_INIT(obj, flags, domain);
}

GRN_API void go_grn_time_init(grn_obj *obj, unsigned char impl_flags) {
	GRN_TIME_INIT(obj, impl_flags);
}

GRN_API void go_grn_time_set(grn_ctx *ctx, grn_obj *obj, long long int unix_usec) {
	GRN_TIME_SET(ctx, obj, unix_usec);
}

GRN_API long long int go_grn_time_value(grn_obj *obj) {
	return GRN_BULK_VSIZE(obj) ? GRN_INT64_VALUE(obj) : 0;
}

char **go_grn_alloc_str_array(int n) {
	return (char **)malloc(sizeof(char *) * n);
}

void go_grn_str_array_set(char **array, int i, char *str) {
	array[i] = str;
}

void go_grn_str_array_free_elems(char **array, int n) {
	for (int i = 0; i < n; i++) {
		free(array[i]);
	}
}

unsigned int *go_grn_alloc_uint_array(int n) {
	return (unsigned int *)malloc(sizeof(unsigned int) * n);
}

void go_grn_uint_array_set(unsigned int *array, int i, unsigned int value) {
	array[i] = value;
}

grn_snip_mapping *go_grn_mapping_html_escape() {
	return GRN_SNIP_MAPPING_HTML_ESCAPE;
}

char *go_grn_malloc_str(int len) {
	return (char *)malloc(sizeof(char) * len);
}
