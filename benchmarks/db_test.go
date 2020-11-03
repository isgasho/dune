package benchmarks

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/scorredoira/dune/lib"
)

func BenchmarkDB_Go(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("CREATE TABLE people (id INTEGER PRIMARY KEY, name TEXT)")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("INSERT INTO people (name) VALUES (?)", "Bob")
		if err != nil {
			log.Fatal(err)
		}

		row := db.QueryRow("SELECT name FROM people LIMIT 1")

		var name string
		if err := row.Scan(&name); err != nil {
			log.Fatal(err)
		}

		if name != "Bob" {
			log.Fatal(name)
		}
	}
}

func BenchmarkDB_Dune(b *testing.B) {
	vm := initVM(b, `
			function main() {
				let db = sql.open("sqlite3", ":memory:")

				db.execRaw("CREATE TABLE people (id INTEGER PRIMARY KEY, name TEXT)")
	
				db.execRaw("INSERT INTO people (name) VALUES (?)", "Bob")
	
				return db.queryValue("SELECT name FROM people LIMIT 1")	
			}
		`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		name, err := vm.Run()
		if err != nil {
			log.Fatal(err)
		}

		if name.ToString() != "Bob" {
			log.Fatal(name)
		}
	}
}

func BenchmarkDB2_Go(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("CREATE TABLE people (id INTEGER PRIMARY KEY, name TEXT)")
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < 10; i++ {
			_, err = db.Exec("INSERT INTO people (name) VALUES (?)", "Bob")
			if err != nil {
				log.Fatal(err)
			}
		}

		for i := 0; i < 100; i++ {
			rows, err := db.Query("SELECT id, name FROM people")
			if err != nil {
				log.Fatal(err)
			}

			var id int
			var name string
			for rows.Next() {
				rows.Scan(&id, &name)
				if id == 0 || name != "Bob" {
					log.Fatal(id, name)
				}
			}
		}
	}
}

func BenchmarkDB2_Dune(b *testing.B) {
	vm := initVM(b, `
			function main() {
				let db = sql.open("sqlite3", ":memory:")
				db.execRaw("CREATE TABLE people (id INTEGER PRIMARY KEY, name TEXT)")
				
				for (let i = 0; i < 10; i++) {
					db.execRaw("INSERT INTO people (name) VALUES (?)", "Bob")
				}
				
				for (let i = 0; i < 100; i++) {
					let rows = db.queryRaw("SELECT id, name FROM people")
					for (let r of rows) {
						if (r.id == 0 || r.name != "Bob") {
							logging.fatal(r.id, r.name)
						}
					}
				}
			}
		`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		if _, err := vm.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
