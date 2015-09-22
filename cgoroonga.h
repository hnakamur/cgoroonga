#ifndef CGOROONGA_H
#define CGOROONGA_H

#include <stdlib.h>
#include <string.h>
#include <groonga/groonga.h>

GRN_API grn_id cgoroonga_obj_header_domain(grn_obj *obj);
GRN_API unsigned char cgoroonga_obj_header_type(grn_obj *obj);
GRN_API unsigned short int cgoroonga_obj_header_flags(grn_obj *obj);

GRN_API void cgoroonga_str_array_set(char **array, int i, char *elem);
GRN_API char *cgoroonga_str_array_get(char **array, int i);
GRN_API void cgoroonga_uint_array_set(unsigned int *array, int i, unsigned int elem);
GRN_API unsigned int cgoroonga_uint_array_get(unsigned int *array, int i);

GRN_API char *cgoroonga_bulk_head(grn_obj *bulk);
GRN_API void cgoroonga_bulk_rewind(grn_obj *bulk);
GRN_API int cgoroonga_bulk_vsize(grn_obj *bulk);

GRN_API long long int cgoroonga_int64_value(grn_obj *obj);

GRN_API grn_snip_mapping *cgoroonga_mapping_html_escape();

GRN_API int cgoroonga_obj_tablep(grn_obj *obj);
GRN_API void cgoroonga_obj_init(grn_obj *obj, unsigned char obj_type, unsigned char impl_flags, grn_id domain);
GRN_API void cgoroonga_value_var_size_init(grn_obj *obj, unsigned char impl_flags, grn_id domain);
GRN_API void cgoroonga_record_init(grn_obj *obj, unsigned char impl_flags, grn_id domain);

GRN_API void cgoroonga_short_text_init(grn_obj *obj, unsigned char impl_flags);
GRN_API void cgoroonga_text_init(grn_obj *obj, unsigned char impl_flags);
GRN_API void cgoroonga_long_text_init(grn_obj *obj, unsigned char impl_flags);
GRN_API void cgoroonga_text_put(grn_ctx *ctx, grn_obj *obj, const char *str, unsigned int len);

GRN_API void cgoroonga_time_init(grn_obj *obj, unsigned char impl_flags);
GRN_API void cgoroonga_time_set(grn_ctx *ctx, grn_obj *obj, long long int unix_usec);

#endif
