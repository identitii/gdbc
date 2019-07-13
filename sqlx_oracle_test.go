// +build cgo,oracle

package gdbc_test

import (
	"fmt"
	"log"
	"testing"

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

func BenchmarkOracleJDBC(b *testing.B) {
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
