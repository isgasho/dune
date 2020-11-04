package sqx

import (
	"testing"
)

func TestAddColumn(t *testing.T) {
	q, err := Parse("alter table foo add bar varchar(6) null")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "ALTER TABLE foo ADD COLUMN bar varchar(6) NULL" {
		t.Fatal(s)
	}
}

func TestAddColumn1(t *testing.T) {
	q, err := Parse("alter table foo add bar mediumblob null")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "ALTER TABLE foo ADD COLUMN bar mediumblob NULL" {
		t.Fatal(s)
	}
}

func TestAddColumn2(t *testing.T) {
	q, err := Parse("alter table foo add column bar varchar(6) null")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "ALTER TABLE foo ADD COLUMN bar varchar(6) NULL" {
		t.Fatal(s)
	}
}

func TestDropColumn(t *testing.T) {
	q, err := Parse("alter table foo drop column bar")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "ALTER TABLE foo DROP COLUMN bar" {
		t.Fatal(s)
	}
}

func TestRenameColumn(t *testing.T) {
	q, err := Parse("alter table foo change fizz bar varchar(6) null")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "ALTER TABLE foo CHANGE fizz bar varchar(6) NULL" {
		t.Fatal(s)
	}
}

func TestModifyColumn(t *testing.T) {
	q, err := Parse("alter table foo modify bar varchar(6) null")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "ALTER TABLE foo MODIFY bar varchar(6) NULL" {
		t.Fatal(s)
	}
}

func TestAddUniqueConstraint(t *testing.T) {
	q, err := Parse("alter table foo add constraint c unique (col1, col2)")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "ALTER TABLE foo ADD CONSTRAINT c UNIQUE (col1, col2)" {
		t.Fatal(s)
	}
}

func TestAddFKConstraint(t *testing.T) {
	q, err := Parse("alter table foo add constraint c foreign key (jj) references bar(id) on delete cascade")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "ALTER TABLE foo ADD CONSTRAINT c FOREIGN KEY(jj) REFERENCES bar(id) ON DELETE CASCADE" {
		t.Fatal(s)
	}
}

func TestDropDatabase(t *testing.T) {
	q, err := Parse("drop database foo")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "DROP DATABASE foo" {
		t.Fatal(s)
	}
}

func TestDropTable(t *testing.T) {
	q, err := Parse("drop table foo")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "DROP TABLE foo" {
		t.Fatal(s)
	}
}

func TestDropTable2(t *testing.T) {
	q, err := Parse("drop table foo.bar")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "DROP TABLE foo.bar" {
		t.Fatal(s)
	}
}

func TestShowTables(t *testing.T) {
	q, err := Parse("show tables")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "SHOW TABLES" {
		t.Fatal(s)
	}
}

func TestShowTables2(t *testing.T) {
	q, err := Parse("show tables from foo")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "SHOW TABLES FROM foo" {
		t.Fatal(s)
	}
}

func TestShowColumns(t *testing.T) {
	q, err := Parse("show columns from foo")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "SHOW COLUMNS FROM foo" {
		t.Fatal(s)
	}
}

func TestShowIndex(t *testing.T) {
	q, err := Parse("show index from foo")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "SHOW INDEX FROM foo" {
		t.Fatal(s)
	}
}

func TestDropIndex(t *testing.T) {
	q, err := Parse("alter table foo drop index bar")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "ALTER TABLE foo DROP INDEX bar" {
		t.Fatal(s)
	}
}

func TestParseCreateDatabaseMysql(t *testing.T) {
	q, err := Parse("create database if not exists foo")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "CREATE DATABASE IF NOT EXISTS foo" {
		t.Fatal(s)
	}
}

func TestParseCreateSqlite(t *testing.T) {
	q, err := Parse("create table if not exists cars (id key, name varchar(10))")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "sqlite3")
	if err != nil {
		t.Fatal(err)
	}

	if s != "CREATE TABLE IF NOT EXISTS cars ("+
		"id INTEGER PRIMARY KEY NOT NULL, "+
		"name VARCHAR(10) NOT NULL COLLATE NOCASE)" {
		t.Fatal(s)
	}
}

func TestParseCreateMysql(t *testing.T) {
	q, err := Parse("create table if not exists cars (id key, name varchar(10))")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "mysql")
	if err != nil {
		t.Fatal(err)
	}

	if s != "CREATE TABLE IF NOT EXISTS cars ("+
		"id int AUTO_INCREMENT NOT NULL, "+
		"name varchar(10) NOT NULL, "+
		"PRIMARY KEY(id))"+
		" ENGINE=InnoDb"+
		" DEFAULT CHARACTER SET = utf8"+
		" DEFAULT COLLATE = utf8_general_ci" {
		t.Fatal(s)
	}
}

func TestParseCreateSqlite1(t *testing.T) {
	q, err := Parse("create table if not exists cars (id key, name varchar(10))")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "sqlite3")
	if err != nil {
		t.Fatal(err)
	}

	if s != "CREATE TABLE IF NOT EXISTS cars ("+
		"id INTEGER PRIMARY KEY NOT NULL, "+
		"name VARCHAR(10) NOT NULL COLLATE NOCASE)" {
		t.Fatal(s)
	}
}

func TestParseCreate2(t *testing.T) {
	q, err := Parse("create table cars (name text, price decimal(12,2) not null)")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "sqlite3")
	if err != nil {
		t.Fatal(err)
	}

	if s != "CREATE TABLE cars (name TEXT NOT NULL COLLATE NOCASE, price REAL(12,2) NOT NULL)" {
		t.Fatal(s)
	}
}

func TestParseCreateUnique(t *testing.T) {
	q, err := Parse("create table cars (code int, price int, constraint code_price unique (code, price))")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "sqlite3")
	if err != nil {
		t.Fatal(err)
	}
	if s != "CREATE TABLE cars (code INTEGER NOT NULL, price INTEGER NOT NULL, CONSTRAINT code_price UNIQUE (code, price))" {
		t.Fatal(s)
	}
}

func TestParseCreate3(t *testing.T) {
	q, err := Parse("create table cars (name text default 'a')")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "sqlite3")
	if err != nil {
		t.Fatal(err)
	}

	if s != "CREATE TABLE cars (name TEXT NOT NULL DEFAULT 'a' COLLATE NOCASE)" {
		t.Fatal(s)
	}
}

func TestParseCreate4(t *testing.T) {
	q, err := Parse("create table users (active bool default true)")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "sqlite3")
	if err != nil {
		t.Fatal(err)
	}

	if s != "CREATE TABLE users (active BOOLEAN NOT NULL DEFAULT true)" {
		t.Fatal(s)
	}
}

func TestParseCreate5(t *testing.T) {
	q, err := Parse("create table users (order int)")
	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "sqlite3")
	if err != nil {
		t.Fatal(err)
	}

	if s != "CREATE TABLE users (order INTEGER NOT NULL)" {
		t.Fatal(s)
	}
}

func TestParseCreateFks(t *testing.T) {
	q, err := Parse(`CREATE TABLE bankaccount (
								id int(11) NOT NULL AUTO_INCREMENT,
								idClient int(11) NOT NULL,
								PRIMARY KEY (id),
								UNIQUE KEY u_name (name),
								CONSTRAINT fk_bankaccountIdClient FOREIGN KEY (idClient) REFERENCES foo_crm_client (id) ON DELETE CASCADE
							)`)

	if err != nil {
		t.Fatal(err)
	}

	s, _, err := toSQL(false, q, "", "sqlite3")
	if err != nil {
		t.Fatal(err)
	}

	if s != "CREATE TABLE bankaccount ("+
		"id INTEGER(11) NOT NULL, idClient INTEGER(11) NOT NULL"+
		", CONSTRAINT u_name UNIQUE (name)"+
		", CONSTRAINT fk_bankaccountIdClient FOREIGN KEY (idClient) REFERENCES foo_crm_client(id))" {
		t.Fatal(s)
	}
}
