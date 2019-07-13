package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lgdbc-driver
#include <stdlib.h>
#include <./libgdbc-driver.h>

graal_isolatethread_t* createIsolate() {
  graal_isolate_t *isolate = NULL;
  graal_isolatethread_t *thread = NULL;

  if (graal_create_isolate(NULL, &isolate, &thread) != 0) {
	return NULL;
  }

  return thread;
}
*/
import "C"

import (
	"log"
	"time"
)

func main() {
	t1 := time.Now()

	// TODO: Release the strings (defer C.free(unsafe.Pointer(cstr)))

	i := C.createIsolate()

	connection := C.newConnection(i, C.CString("org.postgresql.Driver"), C.CString("jdbc:postgresql://localhost/test?loggerLevel=DEBUG"), C.CString("root"), C.CString("password"))
	t2 := time.Now()
	log.Printf("Result from a natively compiled java lib: %v", connection)
	log.Printf("time: %d ms", int64(t2.Sub(t1))/int64(time.Millisecond))

	t1 = time.Now()
	result := C.GoString(C.testQueryJSON(i, connection, C.CString("SELECT 1;")))
	t2 = time.Now()
	log.Printf("time: %d ns", int64(t2.Sub(t1)))

	log.Printf("Result: %s", result)

	// t1 = time.Now()
	// stringResult := C.GoString(C.to_upper(C.createIsolate(), C.CString("elliot")))
	// t2 = time.Now()
	// log.Printf("time: %d ms", int64(t2.Sub(t1))/int64(time.Millisecond))

	// log.Printf("String from java: %s", stringResult)

}
