#ifndef CGOROONGA_H
#define CGOROONGA_H

#include <stdlib.h>
#include <string.h>
#include <groonga/groonga.h>

GRN_API void go_grn_bulk_rewind(grn_obj *bulk);
GRN_API char *go_grn_bulk_head(grn_obj *bulk);
GRN_API int go_grn_bulk_vsize(grn_obj *bulk);

GRN_API grn_obj *go_grn_db_open_or_create(grn_ctx *ctx, const char *path, grn_db_create_optarg *optarg);

GRN_API void go_grn_text_init(grn_obj *text, unsigned char impl_flags);
GRN_API grn_rc go_grn_text_put(grn_ctx *ctx, grn_obj *bulk, const char *str, unsigned int len);
GRN_API void go_grn_record_init(grn_obj *obj, unsigned char flags, grn_id domain);

GRN_API void go_grn_time_init(grn_obj *obj, unsigned char impl_flags);
GRN_API void go_grn_time_set(grn_ctx *ctx, grn_obj *obj, long long int unix_usec);
GRN_API long long int go_grn_time_value(grn_obj *obj);

char **go_grn_alloc_str_array(int n);
void go_grn_str_array_set(char **array, int i, char *str);
void go_grn_str_array_free_elems(char **array, int n);
unsigned int *go_grn_alloc_uint_array(int n);
void go_grn_uint_array_set(unsigned int *array, int i, unsigned int value);
grn_snip_mapping *go_grn_mapping_html_escape();
char *go_grn_malloc_str(int len);
#endif
