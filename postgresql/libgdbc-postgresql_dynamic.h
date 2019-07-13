#ifndef __LIBGDBC_POSTGRESQL_H
#define __LIBGDBC_POSTGRESQL_H

#include <graal_isolate_dynamic.h>


#if defined(__cplusplus)
extern "C" {
#endif

typedef char* (*getError_fn_t)(graal_isolatethread_t*);

typedef void (*enableTracing_fn_t)(graal_isolatethread_t*, int);

typedef char* (*openConnection_fn_t)(graal_isolatethread_t*, char*, char*, char*, int);

typedef char* (*closeConnection_fn_t)(graal_isolatethread_t*);

typedef int (*isValid_fn_t)(graal_isolatethread_t*, int);

typedef char* (*begin_fn_t)(graal_isolatethread_t*);

typedef char* (*commit_fn_t)(graal_isolatethread_t*);

typedef char* (*rollback_fn_t)(graal_isolatethread_t*);

typedef int (*prepare_fn_t)(graal_isolatethread_t*, char*);

typedef char* (*closeStatement_fn_t)(graal_isolatethread_t*, int);

typedef int (*numInput_fn_t)(graal_isolatethread_t*, int);

typedef int (*execute_fn_t)(graal_isolatethread_t*, int);

typedef int (*query_fn_t)(graal_isolatethread_t*, int);

typedef int (*getMoreResults_fn_t)(graal_isolatethread_t*, int);

typedef int (*nextResultSet_fn_t)(graal_isolatethread_t*, int);

typedef char* (*columns_fn_t)(graal_isolatethread_t*, int);

typedef int (*next_fn_t)(graal_isolatethread_t*, int);

typedef char* (*setByte_fn_t)(graal_isolatethread_t*, int, int, char);

typedef char (*getByte_fn_t)(graal_isolatethread_t*, int, int);

typedef char* (*setShort_fn_t)(graal_isolatethread_t*, int, int, short);

typedef short (*getShort_fn_t)(graal_isolatethread_t*, int, int);

typedef char* (*setInt_fn_t)(graal_isolatethread_t*, int, int, int);

typedef int (*getInt_fn_t)(graal_isolatethread_t*, int, int);

typedef char* (*setLong_fn_t)(graal_isolatethread_t*, int, int, long long int);

typedef long long int (*getLong_fn_t)(graal_isolatethread_t*, int, int);

typedef char* (*setFloat_fn_t)(graal_isolatethread_t*, int, int, float);

typedef float (*getFloat_fn_t)(graal_isolatethread_t*, int, int);

typedef char* (*setDouble_fn_t)(graal_isolatethread_t*, int, int, double);

typedef double (*getDouble_fn_t)(graal_isolatethread_t*, int, int);

typedef char* (*getBigDecimal_fn_t)(graal_isolatethread_t*, int, int);

typedef char* (*setString_fn_t)(graal_isolatethread_t*, int, int, char*);

typedef char* (*getString_fn_t)(graal_isolatethread_t*, int, int);

typedef char* (*setTimestamp_fn_t)(graal_isolatethread_t*, int, int, long long int);

typedef long long int (*getTimestamp_fn_t)(graal_isolatethread_t*, int, int);

typedef char* (*setNull_fn_t)(graal_isolatethread_t*, int, int);

typedef char* (*testQueryJSON_fn_t)(graal_isolatethread_t*, char*);

#if defined(__cplusplus)
}
#endif
#endif
