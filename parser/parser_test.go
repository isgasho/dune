package parser

import (
	"strings"
	"testing"

	"github.com/scorredoira/dune/ast"
)

func TestParseFuncDirectives(t *testing.T) {
	a, err := ParseStr(`
		// ::directive1
		function bar() {

		}
	`)

	if err != nil {
		t.Fatal(err)
	}

	// ast.Print(a.File)

	fn, ok := a.File.Stms[0].(*ast.FuncDeclStmt)
	if !ok {
		t.Fatalf("Expected FuncDeclStmt, got %T", a.File.Stms[0])
	}

	if len(fn.Directives) != 1 {
		t.Fatal(fn.Directives)
	}
	if fn.Directives[0] != "directive1" {
		t.Fatal(fn.Directives)
	}
}

func TestParseFuncDirectives2(t *testing.T) {
	a, err := ParseStr(`
		// ::directive ignore

		// ::directive1
		// ::directive2 foo
		function bar() {

		}
	`)

	if err != nil {
		t.Fatal(err)
	}

	// ast.Print(a.File)

	if len(a.File.Directives) != 1 {
		t.Fatal(a.File.Directives)
	}

	fn, ok := a.File.Stms[0].(*ast.FuncDeclStmt)
	if !ok {
		t.Fatalf("Expected FuncDeclStmt, got %T", a.File.Stms[0])
	}

	if len(fn.Directives) != 2 {
		t.Fatal(fn.Directives)
	}
	if fn.Directives[0] != "directive1" || fn.Directives[1] != "directive2 foo" {
		t.Fatal(fn.Directives)
	}
}

func TestParseClassDirectives(t *testing.T) {
	a, err := ParseStr(`
		// ::directive1
		class bar { }
	`)

	if err != nil {
		t.Fatal(err)
	}

	// ast.Print(a.File)

	class, ok := a.File.Stms[0].(*ast.ClassDeclStmt)
	if !ok {
		t.Fatalf("Expected FuncDeclStmt, got %T", a.File.Stms[0])
	}

	if len(class.Directives) != 1 {
		t.Fatal(class.Directives)
	}
	if class.Directives[0] != "directive1" {
		t.Fatal(class.Directives)
	}
}
func TestParseSelector1(t *testing.T) {
	a, err := ParseStr(`let a = b.c.d`)
	if err != nil {
		t.Fatal(err)
	}

	// ast.Print(a.File)

	exp, ok := a.File.Stms[0].(*ast.VarDeclStmt)
	if !ok {
		t.Fatalf("Expected VarDeclStmt, got %T", a.File.Stms[0])
	}

	sel, ok := exp.Value.(*ast.SelectorExpr)
	if !ok {
		t.Fatalf("Expected SelectorExpr, got %T", exp.Value)
	}

	if !sel.First {
		t.Fatal("Not first")
	}
}

func TestParseSelector2(t *testing.T) {
	a, err := ParseStr(`let a = b?.()`)
	if err != nil {
		t.Fatal(err)
	}

	// ast.Print(a.File)

	exp, ok := a.File.Stms[0].(*ast.VarDeclStmt)
	if !ok {
		t.Fatalf("Expected VarDeclStmt, got %T", a.File.Stms[0])
	}

	call, ok := exp.Value.(*ast.CallExpr)
	if !ok {
		t.Fatalf("Expected CallExpr, got %T", exp.Value)
	}

	if !call.First {
		t.Fatal("Not first")
	}
}

func TestParseSelector3(t *testing.T) {
	a, err := ParseStr(`let a = b?.[0]`)
	if err != nil {
		t.Fatal(err)
	}

	//ast.Print(a.File)

	exp, ok := a.File.Stms[0].(*ast.VarDeclStmt)
	if !ok {
		t.Fatalf("Expected VarDeclStmt, got %T", a.File.Stms[0])
	}

	i, ok := exp.Value.(*ast.IndexExpr)
	if !ok {
		t.Fatalf("Expected IndexExpr, got %T", exp.Value)
	}

	if !i.First {
		t.Fatal("Not first")
	}
}

func TestParseSwitchFallthrough1(t *testing.T) {
	_, err := ParseStr(`
		switch(1) {
		case 1:
		case 2:
			let a = 3
		}
		
	`)

	if err != nil {
		t.Fatal(err)
	}
}

func TestParseSwitchFallthrough2(t *testing.T) {
	_, err := ParseStr(`
		switch(1) {
		case 1:
			let a = 3

		case 2:
		}
		
	`)

	if err == nil || !strings.Contains(err.Error(), "Fallthrough") {
		t.Fatal(err)
	}
}

func TestParseSwitchFallthrough3(t *testing.T) {
	_, err := ParseStr(`
		switch(1) {
		case 1:
			let a = 3

		default:

		case 2:
			let b = 3
		}
		
	`)

	if err == nil || !strings.Contains(err.Error(), "Fallthrough") {
		t.Fatal(err)
	}
}
