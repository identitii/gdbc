// +build cgo,oracle

package gdbc_test

import (
	"testing"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/identitii/gdbc/oracle"
	// _ "github.com/mattn/go-oci8"
	// _ "gopkg.in/goracle.v2"
)

var db = env("ORACLE_DB", "xe")
var user = env("ORACLE_USER", "GDBCUSER")
var password = env("ORACLE_PASSWORD", "password")
var host = env("ORACLE_HOST", "localhost")

var jdbcURL = fmt.Sprintf("jdbc:oracle:thin:%s/%s@%s:1521:%s", user, password, host, db)

func aBenchmarkOracleJDBC(b *testing.B) {
	log.Println(jdbcURL)
	benchmark(b, "gdbc-oracle", jdbcURL, sqlx.QUESTION)
}

// FIXME: ORA-12514: TNS:listener does not currently know of service requested in connect descriptor
// func BenchmarkOracleGoracle(b *testing.B) {
// 	benchmark(b, "goracle", fmt.Sprintf("%s/%s@%s/%s", user, password, host, db), sqlx.QUESTION)
// }

// FIXME: Package oci8 was not found in the pkg-config search path.
// func BenchmarkOracleOCI8(b *testing.B) {
// 	benchmark(b, "oci8", fmt.Sprintf("%s/%s@%s/%s", user, password, host, db), sqlx.QUESTION)
// }

func aTestBuildDate(t *testing.T) {

	conn, err := oracle.Open(jdbcURL)
	if err != nil {
		panic(err)
	}

	date, err := conn.DriverBuildDate()
	if err != nil {
		panic(err)
	}

	log.Printf("Oracle driver build date: " + date);

}

func TestContinuousQueryNotification(t *testing.T) {

	oracle.EnableTracing(false)

	

	conn, err := oracle.Open(jdbcURL)
	if err != nil {
		panic(err)
	}

	dcn, err := conn.RegisterDatabaseChangeNotification(1, oracle.CQNOptions{
		NTF_LOCAL_HOST: "192.168.128.187",
		DCN_NOTIFY_ROWIDS: true,
		DCN_QUERY_CHANGE_NOTIFICATION: true,
	})
	if err != nil {
		panic(err)
	}

	state, err := dcn.GetState()
	if err != nil {
		panic(err)
	}
	log.Printf("state: %s", state)

	
	go func() {

		db := getDb(t, "gdbc-oracle", jdbcURL)
		if db == nil {
			return
		}

		for {
			time.Sleep(time.Second * 2)
			log.Println("adding rows")
			runTest(db, sqlx.QUESTION)

			
		}
	}()

	log.Printf("CQN regId: %d", dcn.GetRegId())

	//runTest(db, sqlx.QUESTION)

	state, err = dcn.GetState()
	if err != nil {
		panic(err)
	}
	log.Printf("state: %s", state)

	if err := dcn.AddQuery(`SELECT * FROM "person"`); err != nil {
		panic(err)
	}

	state, err = dcn.GetState()
	if err != nil {
		panic(err)
	}
	log.Printf("state: %s", state)

	// tables, err := dcn.GetTables()
	// if err != nil {
	// 	panic(err)
	// }

	//log.Printf("CQN tables: %#v", tables)

	log.Println("Sleeping")

	time.Sleep(time.Second * 99999)
  
}