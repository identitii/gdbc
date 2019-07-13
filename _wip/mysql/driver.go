package mysql

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lgdbc-mysql
#include <stdlib.h>
#include <./libgdbc-mysql.h>


// TODO: What are isolates? How do they? WHat do? Huh?
graal_isolatethread_t* mysqlCreateIsolate() {

  graal_isolate_t *isolate = NULL;
  graal_isolatethread_t *thread = NULL;

  if (graal_create_isolate(NULL, &isolate, &thread) != 0) {
	return NULL;
  }

  return thread;
}

int mysqlDestroyIsolate(graal_isolatethread_t* thread) {
	return graal_detach_thread(thread);
}
*/
import "C"

import (
	"errors"
	"time"
	"strings"
	"log"

	"github.com/identitii/gdbc"
)

func init() {
	gdbc.Register(&driver{})
}

var tracingEnabled = false

var b = map[bool]C.int{
	false: C.int(0),
	true:  C.int(1),
}

type driver struct {
}

func (d *driver) Open(url, user, password string, txIsolation gdbc.TransactionIsolation) (gdbc.JDBCConnection, error) {
	// TODO: Free the strings! (defer C.free(unsafe.Pointer(cstr)))

	var isolate = C.mysqlCreateIsolate()

	if tracingEnabled {
		C.enableTracing(isolate, b[true])
	}

	C.openConnection(isolate, C.CString(url), C.CString(user), C.CString(password), C.int(txIsolation))

	c := &conn{
		isolate: isolate,
	}

	return c, c.getLastError()
}

// EnableTracing turns on trace level logging in the JDBC wrapper. Note: This is applied at connection time
func (d *driver) EnableTracing(enable bool) {
	tracingEnabled = enable
}

type conn struct {
	isolate *C.graal_isolatethread_t
}

func (c *conn) Close() error {
	C.closeConnection(c.isolate)
	ret := C.mysqlDestroyIsolate(c.isolate)
	log.Printf("destroy isolate retval: %v", ret)
	return nil
}

func (c *conn) Begin() error {
	return err(C.begin(c.isolate))
}

func (c *conn) Commit() error {
	return err(C.commit(c.isolate))
}

func (c *conn) Rollback() error {
	return err(C.rollback(c.isolate)) 
}

func (c *conn) getLastError() error {
	return err(C.getError(c.isolate))
}

func err(ce *C.char) error {
	s := C.GoString(ce)
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func (c *conn) Prepare(sql string) (statement int, err error) {
	return int(C.prepare(c.isolate, C.CString(sql))), c.getLastError()
}

func (c *conn) NumInput(statement int) (inputs int, err error) {
	return int(C.numInput(c.isolate, C.int(statement))), c.getLastError()
}

func (c *conn) Execute(statement int) (updated int, err error) {
	return int(C.execute(c.isolate, C.int(statement))), c.getLastError() 
}

func (c *conn) Query(statement int) (hasResults bool, err error) {
	return int(C.query(c.isolate, C.int(statement))) != 0, c.getLastError() 
}

func (c *conn) Columns(statement int) (columnNames []string, columnTypes []string, err error) {
	columns := C.GoString(C.columns(c.isolate, C.int(statement)))

	err = c.getLastError()
	if err != nil {
		return
	}

	split := strings.Split(columns, "|")

	return strings.Split(split[0], ","), strings.Split(split[1], ","),  nil
}

func (c *conn) Next(statement int) (hasNext bool, err error) {
	return int(C.next(c.isolate, C.int(statement))) != 0, c.getLastError() 
}
 
func (c *conn) SetByte(statement int, index int, value byte) error {
	return err(C.setByte(c.isolate, C.int(statement), C.int(index), C.char(value)))
}

func (c *conn) GetByte(statement int, index int) (byte, error) {
	return byte(C.getByte(c.isolate, C.int(statement), C.int(index))), c.getLastError()
} 

func (c *conn) SetShort(statement int, index int, value int8) error {
	return err(C.setShort(c.isolate, C.int(statement), C.int(index), C.short(value)))
}

func (c *conn) GetShort(statement int, index int) (int8, error) {
	return int8(C.getShort(c.isolate, C.int(statement), C.int(index))), c.getLastError()
}

func (c *conn) SetInt(statement int, index int, value int32) error {
	return err(C.setInt(c.isolate, C.int(statement), C.int(index), C.int(value)))
}

func (c *conn) GetInt(statement int, index int) (int32, error) {
	return int32(C.getInt(c.isolate, C.int(statement), C.int(index))), c.getLastError()
}

func (c *conn) SetLong(statement int, index int, value int64) error {
	return err(C.setLong(c.isolate, C.int(statement), C.int(index), C.longlong(value)))
}

func (c *conn) GetLong(statement int, index int) (int64, error) {
	return int64(C.getLong(c.isolate, C.int(statement), C.int(index))), c.getLastError()
}

func (c *conn) SetFloat(statement int, index int, value float32) error {
	return err(C.setFloat(c.isolate, C.int(statement), C.int(index), C.float(value)))
}

func (c *conn) GetFloat(statement int, index int) (float32, error) {
	return float32(C.getFloat(c.isolate, C.int(statement), C.int(index))), c.getLastError()
}

func (c *conn) SetDouble(statement int, index int, value float64) error {
	return err(C.setDouble(c.isolate, C.int(statement), C.int(index), C.double(value)))
}

func (c *conn) GetDouble(statement int, index int) (float64, error) {
	return float64(C.getDouble(c.isolate, C.int(statement), C.int(index))), c.getLastError()
}

func (c *conn) SetString(statement int, index int, value string) error {
	return err(C.setString(c.isolate, C.int(statement), C.int(index), C.CString(value)))
}

func (c *conn) GetString(statement int, index int) (string, error) {
	return C.GoString(C.getString(c.isolate, C.int(statement), C.int(index))), c.getLastError()
}

func (c *conn) SetTimestamp(statement int, index int, value time.Time) error {
	return err(C.setTimestamp(c.isolate, C.int(statement), C.int(index), C.longlong(value.Unix())))
}

func (c *conn) GetTimestamp(statement int, index int) (time.Time, error) {
	t := int64(C.getTimestamp(c.isolate, C.int(statement), C.int(index)))
	if err := c.getLastError(); err != nil {
		return time.Unix(0, 0), err
	}
	return time.Unix(t, 0), nil
}

func (c *conn) SetNull(statement int, index int) error {
	return err(C.setNull(c.isolate, C.int(statement), C.int(index)))
}

// TestQueryJSON is a quick way to run a query during testing. The resultset is returned as a json array.
func (c *conn) TestQueryJSON(query string) (string, error) {
	result := C.GoString(C.testQueryJSON(c.isolate, C.CString(query)))

	return result, c.getLastError()
}
