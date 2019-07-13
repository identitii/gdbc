// +build cgo,mssql

package gdbc_test

import (
	"net/url"
	"testing"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	//"github.com/identitii/gdbc"

	"github.com/identitii/gdbc/mssql"
	_ "github.com/denisenkom/go-mssqldb"
)

var db = env("MSSQL_DB", "test")
var user = env("MSSQL_USER", "root")
var password = env("MSSQL_PASSWORD", "yourStrong(!)Password")
var host = env("MSSQL_HOST", "localhost")

var jdbcURL = fmt.Sprintf("%s:%s@jdbc:sqlserver://%s:1433;databaseName=%s", user, password, host, db)

func TestPing(t *testing.T) {

	mssql.EnableTracing(false)

	db := getDb(t, "gdbc-mssql", jdbcURL)
	if db == nil {
		return
	}
	defer db.Close()

	log.Printf("Ping: %t", db.Ping())

}

func TestMultipleResultSet(t *testing.T) {
	mssql.EnableTracing(false)

	db := getDb(t, "gdbc-mssql", jdbcURL)
	if db == nil {
		return
	}
	defer db.Close()

	type result1 struct {
		hi int
		there int
	}

	type result2 struct {
		how string
		are string
		you string
	}

	rows, err := db.Query(`SELECT 1 as hi, 2 as there; SELECT '2' as how, '3' as are, '4' as you;`)
	if err != nil {
		panic(err)
	}

	if !rows.Next() {
		panic("result set 1 should have a row")
	}
	
	var r1 result1
	err = rows.Scan(&r1.hi, &r1.there)
	if err != nil {
		panic(err)
	}
	log.Printf("resultset1: %#v", r1)

	if rows.Next() {
		panic("result set 1 should only have one row")
	}
	
	if !rows.NextResultSet() {
		panic("there should have been two result sets")
	}

	if !rows.Next() {
		panic("result set 2 should have a row")
	}

	var r2 result2
	err = rows.Scan(&r2.how, &r2.are, &r2.you)
	if err != nil {
		panic(err)
	}
	log.Printf("resultset2: %#v", r2)

	if rows.Next() {
		panic("result set 2 should only have one row")
	}

}

func BenchmarkMSSQLJDBC(b *testing.B) {
	benchmark(b, "gdbc-mssql", jdbcURL, sqlx.QUESTION)
}

func BenchmarkMSSQLGo(b *testing.B) {
	query := url.Values{}
	query.Add("database", db)
	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%d", host, 1433),
		// Path:  instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}

	benchmark(b, "sqlserver", u.String(), sqlx.DOLLAR)
}
