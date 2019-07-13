//
// Provides an example of the jmoiron/sqlx data mapping library with sqlite
//
package main_test

import (
	"database/sql"
	"fmt"
    "testing"

    "github.com/jmoiron/sqlx"

    "github.com/identitii/gdbc"
    // _ "github.com/identitii/gdbc/postgresql"
    _ "github.com/identitii/gdbc/mssql"
    _ "github.com/lib/pq"
)

var schema = `
DROP TABLE IF EXISTS person CASCADE;
CREATE TABLE person (
    first_name text,
    last_name text,
    email text
);

DROP TABLE IF EXISTS place CASCADE;
CREATE TABLE place (
    country text,
    city text NULL,
    telcode integer
)`

type Person struct {
    FirstName string `db:"first_name"`
    LastName  string `db:"last_name"`
    Email     string
}

type Place struct {
    Country string
    City    sql.NullString
    TelCode int
}

func TestMSSQLJDBC(t *testing.T) {

    db, err := sqlx.Connect("jdbc", "sa:yourStrong(!)Password@jdbc:sqlserver://localhost:1433;databaseName=test")
    if err != nil {
        panic(err)
    }

    runTest(db, sqlx.QUESTION)
    
}

func _BenchmarkPostgresJDBC(b *testing.B) {

    db, err := sqlx.Connect("jdbc", "root:password@jdbc:postgresql://localhost/test?loggerLevel=DEBUG")
    if err != nil {
        panic(err)
    }

    for n := 0; n < b.N; n++ {
        runTest(db, sqlx.QUESTION)
    }
}

func _BenchmarkPostgresGo(b *testing.B) {
    db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "root", "password", "test"))
    if err != nil {
        panic(err)
    }

    for n := 0; n < b.N; n++ {
        runTest(db, sqlx.DOLLAR)
    }
}

func runTest(db *sqlx.DB, bindType int) {
	gdbc.EnableTracing(true)

    
    // exec the schema or fail; multi-statement Exec behavior varies between
    // database drivers;  pq will exec them all, sqlite3 won't, ymmv
    db.MustExec(schema)

    tx := db.MustBegin()
    tx.MustExec(sqlx.Rebind(bindType, "INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)"), "Jason", "Moiron", "jmoiron@jmoiron.net")
    tx.MustExec(sqlx.Rebind(bindType, "INSERT INTO person (first_name, last_name, email) VALUES (?, ?, ?)"), "John", "Doe", "johndoeDNE@gmail.net")
    tx.MustExec(sqlx.Rebind(bindType, "INSERT INTO place (country, city, telcode) VALUES (?, ?, ?)"), "United States", "New York", 1)
    tx.MustExec(sqlx.Rebind(bindType, "INSERT INTO place (country, telcode) VALUES (?, ?)"), "Hong Kong", 852)
    tx.MustExec(sqlx.Rebind(bindType, "INSERT INTO place (country, telcode) VALUES (?, ?)"), "Singapore", 65)
    // Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
    tx.NamedExec("INSERT INTO person (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", &Person{"Jane", "Citizen", "jane.citzen@example.com"})
    tx.Commit()

    // Query the database, storing results in a []Person (wrapped in []interface{})
    people := []Person{}
    err := db.Select(&people, "SELECT * FROM person ORDER BY first_name ASC")
    if err != nil {
        panic(err)
    }
    jason, john := people[0], people[1]

    fmt.Printf("%#v\n%#v", jason, john)
    // Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}
    // Person{FirstName:"John", LastName:"Doe", Email:"johndoeDNE@gmail.net"}

    // You can also get a single result, a la QueryRow
    jason = Person{}
    err = db.Get(&jason, sqlx.Rebind(bindType, "SELECT * FROM person WHERE first_name=?"), "Jason")
    if err != nil {
        panic(err)
    }
    fmt.Printf("%#v\n", jason)
    // Person{FirstName:"Jason", LastName:"Moiron", Email:"jmoiron@jmoiron.net"}

    // if you have null fields and use SELECT *, you must use sql.Null* in your struct
    places := []Place{}
    err = db.Select(&places, "SELECT * FROM place ORDER BY telcode ASC")
    if err != nil {
        panic(err)
    }
    usa, singsing, honkers := places[0], places[1], places[2]
    
    fmt.Printf("%#v\n%#v\n%#v\n", usa, singsing, honkers)
    // Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
    // Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}
    // Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}

    // Loop through rows using only one struct
    place := Place{}
    rows, err := db.Queryx("SELECT * FROM place")
    for rows.Next() {
        err := rows.StructScan(&place)
        if err != nil {
            panic(err)
        } 
        fmt.Printf("%#v\n", place)
    }
    // Place{Country:"United States", City:sql.NullString{String:"New York", Valid:true}, TelCode:1}
    // Place{Country:"Hong Kong", City:sql.NullString{String:"", Valid:false}, TelCode:852}
    // Place{Country:"Singapore", City:sql.NullString{String:"", Valid:false}, TelCode:65}

    // Named queries, using `:name` as the bindvar.  Automatic bindvar support
    // which takes into account the dbtype based on the driverName on sqlx.Open/Connect
    _, err = db.NamedExec(`INSERT INTO person (first_name,last_name,email) VALUES (:first,:last,:email)`, 
        map[string]interface{}{
            "first": "Bin",
            "last": "Smuth",
            "email": "bensmith@allblacks.nz",
    })
    if err != nil {
        panic(err)
    }

    // Selects Mr. Smith from the database
    rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:fn`, map[string]interface{}{"fn": "Bin"})
    if err != nil {
        panic(err)
    }

    // Named queries can also use structs.  Their bind names follow the same rules
    // as the name -> db mapping, so struct fields are lowercased and the `db` tag
    // is taken into consideration.
    rows, err = db.NamedQuery(`SELECT * FROM person WHERE first_name=:first_name`, jason)
    if err != nil {
        panic(err)
    }
}