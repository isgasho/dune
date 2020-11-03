package lib

import (
	"testing"

	"github.com/scorredoira/dune"

	_ "github.com/mattn/go-sqlite3"
)

func TestQuery1(t *testing.T) {
	v := runTest(t, `
		function main() {
			let db = sql.open("sqlite3", ":memory:")
			db.exec("CREATE TABLE foo (id key, name varchar(10))")
			db.exec("INSERT INTO foo VALUES (1, 'a')")
			db.exec("INSERT INTO foo VALUES (2, 'b')")
			return db.queryValue("SELECT count(*) FROM foo WHERE id in ?", [1,2])
		}
	`)

	if v != dune.NewValue(2) {
		t.Fatal(v)
	}
}

func TestSelectAlternateJoin(t *testing.T) {
	v := runTest(t, `
		function main() {
			let db = sql.open("sqlite3", ":memory:")
			db.exec("CREATE TABLE foo (id key, name varchar(10))")
			db.exec("INSERT INTO foo VALUES (1, 'a')")
			db.exec("INSERT INTO foo VALUES (2, 'b')")
			db.exec("CREATE TABLE bar (id key, name varchar(10))")
			db.exec("INSERT INTO bar VALUES (1, 'c')")
			db.exec("INSERT INTO bar VALUES (2, 'a')")
			db.exec("INSERT INTO bar VALUES (3, 'a')")

			let s = "select count(*) from foo f, bar b"
			let q = sql.parse(s)
			q.where("f.name = b.name")

			return db.queryValue(q)
		}
	`)

	if v != dune.NewValue(2) {
		t.Fatal(v)
	}
}

func TestSQLUpdate(t *testing.T) {
	v := runTest(t, `
		function main() {
			let db = sql.open("sqlite3", ":memory:")
			db.exec("CREATE TABLE foo (id key, name varchar(10))")
			db.exec("INSERT INTO foo VALUES (1, 'a')")
			db.exec("INSERT INTO foo VALUES (2, 'b')")
			db.exec("UPDATE foo SET name = 'c' WHERE id = 1")
			return db.queryValue("SELECT name FROM foo WHERE id = 1")
		}
	`)

	if v != dune.NewValue("c") {
		t.Fatal(v)
	}
}

func TestBuilderUpdate(t *testing.T) {
	v := runTest(t, `
		function main() {
			let db = sql.open("sqlite3", ":memory:")
			db.exec("CREATE TABLE foo (id key, name varchar(10))")
			db.exec("INSERT INTO foo VALUES (1, 'a')")
			db.exec("INSERT INTO foo VALUES (2, 'b')")

			let q = sql.parse("UPDATE foo")
			q.addColumns("name = 'c'")
			q.where('id = 1')

			let s = q.toSQL()
			db.exec(s)

			return db.queryValue("SELECT name FROM foo WHERE id = 1")
		}
	`)

	if v != dune.NewValue("c") {
		t.Fatal(v)
	}
}

func TestBuilderDelete(t *testing.T) {
	v := runTest(t, `
		function main() {
			let db = sql.open("sqlite3", ":memory:")
			db.exec("CREATE TABLE foo (id key, name varchar(10))")
			db.exec("INSERT INTO foo VALUES (1, 'a')")
			db.exec("INSERT INTO foo VALUES (2, 'b')")

			let q = sql.parse("DELETE FROM foo")
			q.where('id = 1')

			let s = q.toSQL()

			let rowsAffected = db.exec(s).rowsAffected
			if(rowsAffected != 1) {
				throw "RowsAffected " + rowsAffected
			}
			
			return db.queryValue("SELECT name FROM foo WHERE id = 1")
		}
	`)

	if v.Type != dune.Null {
		t.Fatal(v)
	}
}
