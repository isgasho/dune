package templates

import (
	"strings"
	"testing"

	"github.com/scorredoira/dune"
)

func TestCodeRemovesNewline(t *testing.T) {
	buf, _, err := Compile("<% %>\nFoo")
	if err != nil {
		t.Fatal(err)
	}

	result := string(buf)

	if !strings.Contains(result, "w.write(`Foo`)\n") {
		t.Fatal(result)
	}
}

func TestDirectives(t *testing.T) {
	buf, _, err := Compile(`
	<%@ // ::foo %>
	
	<%@ // ::foo bar %>

	`)
	if err != nil {
		t.Fatal(err)
	}

	result := string(buf)

	p, err := dune.CompileStr(result)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.Directives) != 2 {
		t.Fatal(p.Directives)
	}

	if p.Directives[0] != "foo" {
		t.Fatal(p.Directives[0])
	}

	if p.Directives[1] != "foo bar" {
		t.Fatal(p.Directives[1])
	}
}
