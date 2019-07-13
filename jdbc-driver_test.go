package gdbc_test

// import (
// 	"log"
// 	"testing"

// 	"github.com/identitii/gdbc"
// 	//_ "github.com/identitii/gdbc/postgresql"
// 	//_ "github.com/identitii/gdbc/mssql"
// )

// func aTestPostgresConnection(t *testing.T) {

// 	gdbc.EnableTracing(true)

// 	conn, err := gdbc.Open("root:password@jdbc:postgresql://127.0.0.1:5432/test?loggerLevel=DEBUG", gdbc.TRANSACTION_SERIALIZABLE)
// 	if err != nil {
// 		panic(err)
// 	}

// 	result, err := conn.TestQueryJSON("SELECT 1 as heybuddy, 'how are you' as greeting;")
// 	if err != nil {
// 		panic(err)
// 	}

// 	if result != `[{"greeting":"how are you","heybuddy":1}]` {
// 		panic("bad result: " + result)
// 	}

// 	log.Printf("got result: " + result)

// 	if err := conn.Close(); err != nil {
// 		panic(err)
// 	}
// }

// func aTestMSSQLConnection(t *testing.T) {

// 	gdbc.EnableTracing(true)

// 	conn, err := gdbc.Open("sa:yourStrong(!)Password@jdbc:sqlserver://localhost:1433;databaseName=test", gdbc.TRANSACTION_SERIALIZABLE)
// 	if err != nil {
// 		panic(err)
// 	}

// 	result, err := conn.TestQueryJSON("SELECT 1 as heybuddy, 'how are you' as greeting;")
// 	if err != nil {
// 		panic(err)
// 	}

// 	if result != `[{"greeting":"how are you","heybuddy":1}]` {
// 		panic("bad result: " + result)
// 	}

// 	log.Printf("got result: " + result)

// 	if err := conn.Close(); err != nil {
// 		panic(err)
// 	}
// }

// func _TestOracleConnection(t *testing.T) {

// 	gdbc.EnableTracing(true)

// 	conn, err := gdbc.Open("sys as sysdba:Oradoc_db1@jdbc:oracle:thin:sys as sysdba/Oradoc_db1@localhost:1521:ORCLCDB", gdbc.TRANSACTION_READ_COMMITTED)
// 	if err != nil {
// 		panic(err)
// 	}

// 	result, err := conn.TestQueryJSON(`SELECT 1 AS heybuddy, 'how are you' AS greeting FROM DUAL`)
// 	if err != nil {
// 		panic(err)
// 	}

// 	if result != `[{"HEYBUDDY":1,"GREETING":"how are you"}]` {
// 		panic("bad result: " + result)
// 	}

// 	log.Printf("got result: " + result)

// 	if err := conn.Close(); err != nil {
// 		panic(err)
// 	}
// }