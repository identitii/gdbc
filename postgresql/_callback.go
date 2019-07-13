package postgresql

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lgdbc-postgresql
#include <./libgdbc-postgresql.h>

*/
import "C"
import "log"

//export on_update
func on_update(t *C.graal_isolatethread_t, name, value *C.char) bool {
	log.Printf("Got an update pushed from java: name: %s value: %s", C.GoString(name), C.GoString(value))
	return true
}
