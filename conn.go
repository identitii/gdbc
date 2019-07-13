package gdbc

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"
)

func ParseJDBCURL(url string) (user, password, u string, err error) {
	var re = regexp.MustCompile(`(?m)(?P<user>.*):(?P<password>.*)@(?P<url>jdbc:.*)`)

	matches := re.FindAllStringSubmatch(url, 1)
	if len(matches) > 0 {
		// URL has user and password at the beginning
		return matches[0][1], matches[0][2], matches[0][3], nil
	}

	// No user and password, just pass the url back
	return "", "", url, nil
}

func NewConn(jdbcConn JDBCConnection) Conn {
	return &conn{
		c: jdbcConn,
	}
}

type Conn interface {
	driver.Conn
	driver.Pinger
	gdbc()
}

type conn struct {
	c  JDBCConnection
	tx *tx
}

func (c *conn) gdbc() {}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	id, err := c.c.Prepare(query)
	if err != nil {
		return nil, err
	}

	return &Stmt{
		c:  c.c,
		id: id,
	}, nil
}

func (c *conn) Close() error {
	return c.c.Close(true)
}

func (c *conn) Begin() (driver.Tx, error) {
	if c.tx != nil {
		return nil, errors.New("transaction has already begun")
	}
	if err := c.c.Begin(); err != nil {
		return nil, err
	}
	c.tx = &tx{
		c: c,
	}
	return c.tx, nil
}

func (c *conn) Ping(ctx context.Context) error {
	valid, err := c.c.IsValid(0)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("connection not valid")
	}
	return nil
}

type Stmt struct {
	c  JDBCConnection
	id int
}

func (s *Stmt) ID() int {
	return s.id
}

func (s *Stmt) Close() error {
	return s.c.CloseStatement(s.id)
}

func (s *Stmt) NumInput() int {
	i, err := s.c.NumInput(s.id)
	if err != nil {
		return -1
	}
	return i
}

func (s *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	if err := s.sendArgs(args); err != nil {
		return nil, err
	}
	updated, err := s.c.Execute(s.id)
	return &result{
		updated: updated,
	}, err
}

func (s *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	if err := s.sendArgs(args); err != nil {
		return nil, err
	}
	hasResults, err := s.c.Query(s.id)
	if err != nil {
		return nil, err
	}
	if !hasResults {
		none := []string{}
		return &rows{
			s:           s,
			columnNames: &none,
		}, nil
	}
	return &rows{
		s: s,
	}, nil
}

func (s *Stmt) sendArgs(args []driver.Value) (err error) {
	for i, a := range args {
		idx := i + 1

		switch arg := a.(type) {
		case byte:
			err = s.c.SetByte(s.id, idx, arg)
		case int8:
			err = s.c.SetShort(s.id, idx, arg)
		case int32:
			err = s.c.SetInt(s.id, idx, arg)
		case int64:
			err = s.c.SetLong(s.id, idx, arg)
		case float32:
			err = s.c.SetFloat(s.id, idx, arg)
		case float64:
			err = s.c.SetDouble(s.id, idx, arg)
		case string:
			err = s.c.SetString(s.id, idx, arg)
		case time.Time:
			err = s.c.SetTimestamp(s.id, idx, arg)
		case nil:
			err = s.c.SetNull(s.id, idx)
		default:
			err = fmt.Errorf("Unhandled param type(%d): %T %v", i, a, a)
		}
		if err != nil {
			return
		}
	}
	return
}

type result struct {
	updated int
}

func (r *result) LastInsertId() (int64, error) {
	return -1, errors.New("Result.LastInsertId not supported")
}

func (r *result) RowsAffected() (int64, error) {
	return int64(r.updated), nil
}

type tx struct {
	c        *conn
	finished bool
}

func (t *tx) finish() bool {
	if t.finished {
		return false
	}
	t.finished = true
	t.c.tx = nil
	return true
}

func (t *tx) Commit() error {
	if !t.finish() {
		return nil
	}

	return t.c.c.Commit()
}

func (t *tx) Rollback() error {
	if !t.finish() {
		return nil
	}

	return t.c.c.Rollback()
}

type Rows interface {
	driver.RowsColumnTypeDatabaseTypeName
	// driver.RowsColumnTypeLength
}

type rows struct {
	s           *Stmt
	columnNames *[]string
	columnTypes *[]string
}

func (r *rows) Columns() []string {

	columnNames, columnTypes, err := r.s.c.Columns(r.s.id)
	if err != nil {
		log.Printf("ERROR: failed to get column names in row: %s", err)
		return nil
	}
	r.columnNames = &columnNames
	r.columnTypes = &columnTypes // TODO: Add sql types so RowsColumnTypeDatabaseTypeName etc can be implemented

	return *r.columnNames
}

func (r *rows) Close() error {
	return r.s.Close()
}

func (r *rows) Next(dest []driver.Value) error {

	hasNext, err := r.s.c.Next(r.s.id)
	if err != nil {
		return err
	}
	if !hasNext {
		return io.EOF
	}

	for i, columnType := range *r.columnTypes {
		var idx = i + 1
		switch columnType { // TODO: Test all these column types
		case "java.lang.Byte":
			dest[i], err = r.s.c.GetByte(r.s.id, idx)
		case "java.lang.Short":
			dest[i], err = r.s.c.GetShort(r.s.id, idx)
		case "java.lang.Integer":
			dest[i], err = r.s.c.GetInt(r.s.id, idx)
		case "java.lang.Long":
			dest[i], err = r.s.c.GetLong(r.s.id, idx)
		case "java.lang.Float":
			dest[i], err = r.s.c.GetFloat(r.s.id, idx)
		case "java.lang.Double":
			dest[i], err = r.s.c.GetDouble(r.s.id, idx)
		case "java.math.BigDecimal":
			dest[i], err = r.s.c.GetBigDecimal(r.s.id, idx)
		case "java.lang.String":
			dest[i], err = r.s.c.GetString(r.s.id, idx)
		case "java.sql.Timestamp":
			dest[i], err = r.s.c.GetTimestamp(r.s.id, idx)
		default:
			err = errors.New("column type " + columnType + " not supported")
		}
		if err != nil {
			return fmt.Errorf("failed reading column %d (%s): %s", idx, (*r.columnNames)[i], err)
		}
	}
	return nil
}

func (r *rows) HasNextResultSet() bool {
	return r.s.c.GetMoreResults(r.s.id) // TODO: This actually goes to the next result set... but I *think* that's ok?
}

func (r *rows) NextResultSet() error {
	if r.s.c.NextResultSet(r.s.id) {
		return nil
	}

	return io.EOF
}
