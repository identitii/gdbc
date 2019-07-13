
// +build cgo,postgresql

package gdbc_test

import (
	"testing"
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/identitii/gdbc/postgresql"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
)

var db = env("POSTGRES_DB", "test")
var user = env("POSTGRES_USER", "root")
var password = env("POSTGRES_PASSWORD", "password")
var host = env("POSTGRES_HOST", "localhost")

func BenchmarkPostgresJDBC(b *testing.B) {
	benchmark(b, "gdbc-postgresql", fmt.Sprintf("%s:%s@jdbc:postgresql://%s/%s?loggerLevel=DEBUG", user, password, host, db), sqlx.QUESTION)
}


func BenchmarkPostgresGoPQ(b *testing.B) {
	benchmark(b, "postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, 5432, user, password, db), sqlx.DOLLAR)
}

func BenchmarkPostgresGoPGX(b *testing.B) {
	benchmark(b, "pgx", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, 5432, user, password, db), sqlx.DOLLAR)
}
