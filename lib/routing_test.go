package lib

import (
	"fmt"
	"testing"
)

func TestRoute(t *testing.T) {
	var tests = []struct {
		route  *httpRoute
		url    string
		match  bool
		params string
	}{
		{NewRoute("*"), "/", true, "*"},
		{NewRoute("*"), "/llll", true, "*"},
		{NewRoute("/foo"), "/foo", true, "*"},
		{NewRoute("/foo/*"), "/foo", true, "*"},
		{NewRoute("/foo/bar/*"), "/foo/bar/cxxx", true, "*"},
	}

	router := newRouter()

	for i, test := range tests {
		router.Reset()

		if test.route.URL != "" {
			router.Add(test.route)
		}

		err := assertRoute(router, test.route, test.url, test.match, test.params)
		if err != nil {
			t.Fatalf("Test %d: %s -> %s", i, test.url, err.Error())
		}
	}
}

func TestRoute1(t *testing.T) {
	var tests = []struct {
		route  *httpRoute
		url    string
		match  bool
		params string
	}{
		{NewRoute("*"), "/", true, "*"},
		{NewRoute("/:foo/bar"), "/foo/bar/xxx", true, "*"},
	}

	router := newRouter()

	for i, test := range tests {
		if test.route.URL != "" {
			router.Add(test.route)
		}

		err := assertRoute(router, test.route, test.url, test.match, test.params)
		if err != nil {
			t.Fatalf("Test %d: %s -> %s", i, test.url, err.Error())
		}
	}
}

func TestRoute2(t *testing.T) {
	var tests = []struct {
		route  *httpRoute
		url    string
		match  bool
		params string
	}{
		{NewRoute("/:foo/*"), "/foo", true, "*"},
		{NewRoute("/:foo/bar.json"), "/foo/bar", true, "*"},
	}

	router := newRouter()

	for i, test := range tests {
		if test.route.URL != "" {
			router.Add(test.route)
		}

		err := assertRoute(router, test.route, test.url, test.match, test.params)
		if err != nil {
			t.Fatalf("Test %d: %s -> %s", i, test.url, err.Error())
		}
	}
}

func TestRoute3(t *testing.T) {
	router := newRouter()
	router.Reset()

	var tests = []struct {
		route  *httpRoute
		url    string
		match  bool
		params string
	}{
		{
			NewRoute("/:namespace/:action"),
			"/app1/test", true, "namespace=app1,action=test",
		},
		{
			NewRoute("/customers"),
			"/customers", true, "",
		},
		{
			// capture action as a parameter
			NewRoute("/customers/:id"),
			"/customers/234", true, "id=234",
		},
		{
			// capture the extension as a parameter
			NewRoute("/foo.:ext"),
			"/foo.js", true, "ext=js",
		},
		{
			NewRoute("/customers/foo"),
			"/customers/foo", true, "",
		},
		{
			NewRoute("/customers/:id/print"),
			"/customers/234/print", true, "id=234",
		},
		{
			NewRoute("/public/*"),
			"/public/images/pics/dog.jpeg", true, "",
		},
		{
			NewRoute("/:namespace/public/*"),
			"/demo/public/images/pics/dog.jpeg", true, "namespace=demo",
		},
		{
			NewRoute(""), "/foo", false, "",
		},
		{
			NewRoute(""), "/customers/10/foo", false, "",
		},
	}

	for _, test := range tests {
		if test.route.URL != "" {
			router.Add(test.route)
		}
	}

	for i, test := range tests {
		err := assertRoute(router, test.route, test.url, test.match, test.params)
		if err != nil {
			t.Fatalf("Test %d: %s -> %s", i+1, test.url, err.Error())
		}
	}
}

func assertRoute(router *httpRouter, r *httpRoute, url string, match bool, params string) error {
	m, ok := router.Match(url)
	if ok != match {
		if ok {
			fmt.Println(m.Route.URL)
		}
		return fmt.Errorf("expected match: %t, got: %t", match, ok)
	}

	if params != "*" {
		items := Split(params, ",")
		if len(items) != len(m.Params) {
			return fmt.Errorf("params: expected %d, got %d", len(items), len(m.Params))
		}

		for _, item := range items {
			parts := Split(item, "=")
			key := parts[0]
			value := parts[1]
			if m.GetParam(key) != value {
				return fmt.Errorf("param %s: expected '%s', got '%s'", key, value, m.GetParam(key))
			}
		}
	}

	return nil
}
