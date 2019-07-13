package gdbc

import (
	"errors"
	"time"
)

// UnsupportedJDBCURL means that no jdbc driver has been registered that accepts the url provided
var UnsupportedJDBCURL = errors.New("unsupported jdbc url")

// TransactionIsolation is passed through to java.sql.Connection.setTransactionIsolation method (
// For more information see java.sql.Connection.TRANSACTION_READ_UNCOMMITTED etc.
type TransactionIsolation int

const (
	TRANSACTION_NONE             TransactionIsolation = 0
	TRANSACTION_READ_UNCOMMITTED TransactionIsolation = 1
	TRANSACTION_READ_COMMITTED   TransactionIsolation = 2
	TRANSACTION_REPEATABLE_READ  TransactionIsolation = 4
	TRANSACTION_SERIALIZABLE     TransactionIsolation = 8
)

// var regoLock sync.Mutex
// var registeredDrivers []JDBCDriver

// func Register(driver JDBCDriver) {
// 	regoLock.Lock()
// 	defer regoLock.Unlock()

// 	if len(registeredDrivers) > 0 {
// 		panic("a jdbc driver has already been registered. only one driver can be imported at a time.") // TODO: FIXME:
// 	}

// 	registeredDrivers = append(registeredDrivers, driver)
// }

// func EnableTracing(enable bool) {
// 	for _, driver := range registeredDrivers {
// 		driver.EnableTracing(enable)
// 	}
// }

// func Open(url string, txIsolation TransactionIsolation) (JDBCConnection, error) {
// 	regoLock.Lock()
// 	defer regoLock.Unlock()

// 	user, password, jdbcURL, err := ParseJDBCURL(url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, driver := range registeredDrivers {
// 		conn, err := driver.Open(jdbcURL, user, password, txIsolation)
// 		if err == nil {
// 			return conn, nil
// 		}
// 		log.Printf("Err message: %s", err)
// 		if err != UnsupportedJDBCURL && !strings.Contains(err.Error(), "No suitable driver found") {
// 			log.Printf("bad error")
// 			return nil, err
// 		}
// 	}
// 	return nil, UnsupportedJDBCURL
// }

type JDBCDriver interface {
	EnableTracing(enable bool)
	Open(url, user, password string, txIsolation TransactionIsolation) (JDBCConnection, error)
}

type JDBCConnection interface {
	Close(keepIsolate bool) error
	Begin() error
	Commit() error
	Rollback() error
	IsValid(timeout int) (valid bool, err error)

	Prepare(sql string) (statement int, err error)
	CloseStatement(statement int) (err error)
	NumInput(statement int) (inputs int, err error)
	Execute(statement int) (updated int, err error)
	Query(statement int) (hasResults bool, err error)
	Columns(statement int) (columnNames []string, columnTypes []string, err error)
	Next(statement int) (hasNext bool, err error)
	GetMoreResults(statement int) bool
	NextResultSet(statement int) bool

	SetByte(statement int, index int, value byte) error
	GetByte(statement int, index int) (byte, error)
	SetShort(statement int, index int, value int8) error
	GetShort(statement int, index int) (int8, error)
	SetInt(statement int, index int, value int32) error
	GetInt(statement int, index int) (value int32, err error)
	SetLong(statement int, index int, value int64) error
	GetLong(statement int, index int) (int64, error)
	SetFloat(statement int, index int, value float32) error
	GetFloat(statement int, index int) (float32, error)
	SetDouble(statement int, index int, value float64) error
	GetDouble(statement int, index int) (float64, error)
	GetBigDecimal(statement int, index int) (float64, error)
	SetString(statement int, index int, value string) error
	GetString(statement int, index int) (value string, err error)
	SetTimestamp(statement int, index int, value time.Time) error
	GetTimestamp(statement int, index int) (time.Time, error)
	SetNull(statement int, index int) error

	TestQueryJSON(query string) (result string, err error)
}
