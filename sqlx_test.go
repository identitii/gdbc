//
// Provides an example of the jmoiron/sqlx data mapping library with sqlite
//
package gdbc_test

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	//_ "github.com/identitii/gdbc/mysql"
	//_ "github.com/identitii/gdbc/oracle"
	// _ "github.com/go-sql-driver/mysql"
	//_ "github.com/mattn/go-oci8"
)

func env(name string, def string) string {
	if os.Getenv(name) != "" {
		return os.Getenv(name)
	}
	return def
}

// 13796900 ns/op - github.com/identitii/gdbc/postgresql
// 21256796 ns/op - github.com/lib/pq
// 24224562 ns/op - github.com/jackc/pgx/stdlib
// 35% and 43% faster

// 19064821 ns/op - github.com/identitii/gdbc/mssql
// 21760250 ns/op - github.com/denisenkom/go-mssqldb
// 12% faster

// 37477937 ns/op - github.com/identitii/gdbc/oracle

// func BenchmarkMySQLJDBC(b *testing.B) {
// 	benchmark(b, "jdbc", "root:password@jdbc:mysql://localhost/test", sqlx.QUESTION)
// }

// func BenchmarkMySQLGo(b *testing.B) {
// 	benchmark(b, "mysql", "root:password@/test", sqlx.QUESTION)
// }

type test interface {
	SkipNow()
	Fail()
}

func benchmark(b *testing.B, name, url string, bindType int) {
	// gdbc.EnableTracing(true)

	// log.Printf("connecting to %s %s", name, url)

	db := getDb(b, name, url)
	if db == nil {
		return
	}
	defer db.Close()

	for n := 0; n < b.N; n++ {
		runTest(db, bindType)
	}
}

func getDb(b test, name, url string) *sqlx.DB {
	db, err := sqlx.Connect(name, url)
	if err != nil {
		if err.Error() == "unsupported jdbc url" {
			// is ok. just don't have this driver
			b.SkipNow()
			return nil
		}

		if strings.Contains(err.Error(), "connection refused") {
			// is ok. server isn't running
			b.SkipNow()
			return nil
		}
		log.Printf("Connection error: %s", err)
		b.Fail()
		return nil
	}

	db.Exec(`DROP TABLE "person"`) // Ok to ignore. it might not be there
	db.Exec(`DROP TABLE "place"`)

	for _, s := range schema {
		db.MustExec(s)
	}
	return db
}

var schema = []string{
	// `DROP TABLE "person"`,
	`CREATE TABLE "person" (
        "first_name" varchar(100),
        "last_name" varchar(100),
        "email" varchar(100)
    )`,

	`CREATE TABLE "place" (
        "country" varchar(100),
        "city" varchar(100) NULL,
        "telcode" integer
    )`,
}

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

type Place struct {
	Country string         `db:"country"`
	City    sql.NullString `db:"city"`
	TelCode int            `db:"telcode"`
}

func runTest(db *sqlx.DB, bindType int) {

	tx := db.MustBegin()
	tx.MustExec(sqlx.Rebind(bindType, `INSERT INTO "person" ("first_name", "last_name", "email") VALUES (?, ?, ?)`), "Jason", "Moiron", "jmoiron@jmoiron.net")
	tx.MustExec(sqlx.Rebind(bindType, `INSERT INTO "person" ("first_name", "last_name", "email") VALUES (?, ?, ?)`), "John", "Doe", "johndoeDNE@gmail.net")
	tx.MustExec(sqlx.Rebind(bindType, `INSERT INTO "place" ("country", "city", "telcode") VALUES (?, ?, ?)`), "United States", "New York", 1)
	tx.MustExec(sqlx.Rebind(bindType, `INSERT INTO "place" ("country", "telcode") VALUES (?, ?)`), "Hong Kong", 852)
	tx.MustExec(sqlx.Rebind(bindType, `INSERT INTO "place" ("country", "telcode") VALUES (?, ?)`), "Singapore", 65)
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	tx.NamedExec(`INSERT INTO "person" ("first_name", "last_name", "email") VALUES (:first_name, :last_name, :email)`, &Person{"Jane", "Citizen", "jane.citzen@example.com"})
	tx.Commit()

	// Query the database, storing results in a []Person (wrapped in []interface{})
	people := []Person{}
	err := db.Select(&people, `SELECT * FROM "person"`)
	if err != nil {
		panic(err)
	}
	//jason, john := people[0], people[1]

	// fmt.Printf("%#v\n%#v", jason, john)
	// Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}
	// Person{FirstName:"John", LastName:"Doe", Email:"johndoeDNE@gmail.net"}

	// You can also get a single result, a la QueryRow
	jason := Person{}
	err = db.Get(&jason, sqlx.Rebind(bindType, `SELECT * FROM "person" WHERE "first_name"=?`), "Jason")
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%#v\n", jason)
	// Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}

	// if you have null fields and use SELECT *, you must use sql.Null* in your struct
	places := []Place{}
	err = db.Select(&places, `SELECT * FROM "place" ORDER BY "telcode" ASC`)
	if err != nil {
		panic(err)
	}
	//usa, singsing, honkers := places[0], places[1], places[2]

	//NOfmt.Printf("%#v\n%#v\n%#v\n", usa, singsing, honkers)
	// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
	// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}
	// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}

	// Loop through rows using only one struct
	place := Place{}
	rows, err := db.Queryx(`SELECT * FROM "place"`)
	for rows.Next() {
		err := rows.StructScan(&place)
		if err != nil {
			panic(err)
		}
		//NOfmt.Printf("%#v\n", place)
	}
	rows.Close()

	// Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
	// Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}
	// Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}

	// Named queries, using `:name` as the bindvar.  Automatic bindvar support
	// which takes into account the dbtype based on the driverName on sqlx.Open/Connect
	_, err = db.NamedExec(`INSERT INTO "person" ("first_name", "last_name", "email") VALUES (:first,:last,:email)`,
		map[string]interface{}{
			"first": "Bin",
			"last":  "Smuth",
			"email": "bensmith@allblacks.nz",
		})
	if err != nil {
		panic(err)
	}

	// Selects Mr. Smith from the database
	rows, err = db.NamedQuery(`SELECT * FROM "person" WHERE "first_name"=:fn`, map[string]interface{}{"fn": "Bin"})
	if err != nil {
		panic(err)
	}
	rows.Close()

	// Named queries can also use structs.  Their bind names follow the same rules
	// as the name -> db mapping, so struct fields are lowercased and the `db` tag
	// is taken into consideration.
	rows, err = db.NamedQuery(`SELECT * FROM "person" WHERE "first_name"=:first_name`, jason)
	if err != nil {
		panic(err)
	}
	rows.Close()
}
