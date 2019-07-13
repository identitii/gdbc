package oracle

/*
#include <./libgdbc-oracle.h>
#include <stdlib.h>
*/
import "C"

import (
	"strings"
	"log"
	"unsafe"
	"encoding/json"
	"database/sql/driver"

	"github.com/identitii/gdbc"
)

var registrations = map[int64]*DatabaseChangeRegistration{}

// CQNOptions are passed through to the java driver
// From https://docs.oracle.com/database/121/JJDBC/dbchgnf.htm#JJDBC28818
type CQNOptions struct  {

	// If set to true, DELETE operations will not generate any database change event.
	DCN_IGNORE_DELETEOP bool `json:",omitempty"`

	// If set to true, INSERT operations will not generate any database change event.
	DCN_IGNORE_INSERTOP bool `json:",omitempty"`

	// If set to true, UPDATE operations will not generate any database change event.
	DCN_IGNORE_UPDATEOP bool `json:",omitempty"`

	// Specifies the number of transactions by which the client is willing to lag behind.
	// Note: If this option is set to any value other than 0, then ROWID level granularity of information will not be available in the events, even if the DCN_NOTIFY_ROWIDS option is set to true.
	DCN_NOTIFY_CHANGELAG int `json:",omitempty"`

	// Database change events will include row-level details, such as operation type and ROWID.
	DCN_NOTIFY_ROWIDS bool `json:",omitempty"`

	// Activates query change notification instead of object change notification.
	// Note: This option is available only when running against an 11.0 database.
	DCN_QUERY_CHANGE_NOTIFICATION bool `json:",omitempty"`

	// Specifies the IP address of the computer that will receive the notifications from the server.
	NTF_LOCAL_HOST string `json:",omitempty"`

	// Specifies the TCP port that the driver should use for the listener socket.
	NTF_LOCAL_TCP_PORT int `json:",omitempty"`

	// Specifies if the registration should be expunged on the first notification event.
	NTF_QOS_PURGE_ON_NTFN bool `json:",omitempty"`

	// Specifies whether or not to make the notifications persistent, which comes at a performance cost.
	NTF_QOS_RELIABLE bool `json:",omitempty"`

	// Specifies the time in seconds after which the registration will be automatically expunged by the database.
	NTF_TIMEOUT int `json:",omitempty"`
}

func (c *Conn) RegisterDatabaseChangeNotification(channelSize int, options CQNOptions) (*DatabaseChangeRegistration, error) {

	optionsJSON, _ := json.Marshal(options)

	coptionsJSON := C.CString(string(optionsJSON))
	defer func() {
		C.free(unsafe.Pointer(coptionsJSON))
	}()

	id := int(C.oracleRegisterDatabaseChangeNotification(c.c.isolate, coptionsJSON))
	err := c.c.getLastError()
	if err != nil {
		return nil, err 
	}

	r := &DatabaseChangeRegistration{
		id: id,
		regID: int64(C.oracleRegistrationGetRegId(c.c.isolate, C.int(id))),
		c: c,
		events: make(chan<-*DatabaseChangeEvent, channelSize),
	}
	if err := c.c.getLastError(); err != nil {
		return nil, err
	}
 
	return r, nil
}

type DatabaseChangeRegistration struct {
	id int
	regID int64

	c *Conn
	events chan<-*DatabaseChangeEvent
}

func (r *DatabaseChangeRegistration) GetRegId() int64 {
	return r.regID
}

func (r *DatabaseChangeRegistration) GetTables() ([]string, error) {
	tables := C.GoString(C.oracleRegistrationGetTables(r.c.c.isolate, C.int(r.id)))
	return strings.Split(tables, ","), r.c.c.getLastError()
}

func (r *DatabaseChangeRegistration) GetState() (string, error) {
	state := C.GoString(C.oracleRegistrationGetState(r.c.c.isolate, C.int(r.id)))
	return state, r.c.c.getLastError()
}

func (r *DatabaseChangeRegistration) AddQuery(query string, args ...driver.Value) error {
	stmt, err := r.c.Prepare(query)
	if err != nil {
		return err
	}

	err = goerr(C.oracleSetDatabaseChangeRegistration(r.c.c.isolate, C.int(stmt.(*gdbc.Stmt).ID()), C.int(r.id)))
	if err != nil {
		return err
	}

	rows, err := stmt.Query(args) 
	if err != nil {
		return err
	}
	log.Printf("columns: %v", rows.Columns())

	// var a,b,c string

	// for {
	// 	log.Printf("hhh")
	// 	if err := rows.Next([]driver.Value{&a, &b, &c}); err != nil {
	// 		log.Printf("row next %s", err)
	// 		break
	// 	}
	// 	log.Printf("%s %s %s", a, b, c)
	// }

	//rows.Close()
	//stmt.Close()

	log.Printf("executed query... %#v", rows)

	return nil
}

type DatabaseChangeEvent struct {
} 

//export oracle_on_change_event
func oracle_on_change_event(regID int64, event *C.char) bool {
	log.Printf("GOT AN ORACLE EVENT! regId: %d event: %s", regID, C.GoString(event))
	return false
}
