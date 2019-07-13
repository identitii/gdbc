#ifndef __LIBGDBC_ORACLE_H
#define __LIBGDBC_ORACLE_H

#include <graal_isolate.h>


#if defined(__cplusplus)
extern "C" {
#endif

char* getError(graal_isolatethread_t*);

void enableTracing(graal_isolatethread_t*, int);

char* openConnection(graal_isolatethread_t*, char*, char*, char*, int);

char* closeConnection(graal_isolatethread_t*);

int isValid(graal_isolatethread_t*, int);

char* begin(graal_isolatethread_t*);

char* commit(graal_isolatethread_t*);

char* rollback(graal_isolatethread_t*);

int prepare(graal_isolatethread_t*, char*);

char* closeStatement(graal_isolatethread_t*, int);

int numInput(graal_isolatethread_t*, int);

int execute(graal_isolatethread_t*, int);

int query(graal_isolatethread_t*, int);

int getMoreResults(graal_isolatethread_t*, int);

int nextResultSet(graal_isolatethread_t*, int);

char* columns(graal_isolatethread_t*, int);

int next(graal_isolatethread_t*, int);

char* setByte(graal_isolatethread_t*, int, int, char);

char getByte(graal_isolatethread_t*, int, int);

char* setShort(graal_isolatethread_t*, int, int, short);

short getShort(graal_isolatethread_t*, int, int);

char* setInt(graal_isolatethread_t*, int, int, int);

int getInt(graal_isolatethread_t*, int, int);

char* setLong(graal_isolatethread_t*, int, int, long long int);

long long int getLong(graal_isolatethread_t*, int, int);

char* setFloat(graal_isolatethread_t*, int, int, float);

float getFloat(graal_isolatethread_t*, int, int);

char* setDouble(graal_isolatethread_t*, int, int, double);

double getDouble(graal_isolatethread_t*, int, int);

char* getBigDecimal(graal_isolatethread_t*, int, int);

char* setString(graal_isolatethread_t*, int, int, char*);

char* getString(graal_isolatethread_t*, int, int);

char* setTimestamp(graal_isolatethread_t*, int, int, long long int);

long long int getTimestamp(graal_isolatethread_t*, int, int);

char* setNull(graal_isolatethread_t*, int, int);

char* testQueryJSON(graal_isolatethread_t*, char*);

#if defined(__cplusplus)
}
#endif
#endif
