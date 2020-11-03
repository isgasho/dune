package benchmarks

import (
	"log"
	"testing"
)

func BenchmarkClassField(b *testing.B) {
	vm := initVM(b, `
			let foo = new Foo()

			class Foo {
				a = 3
			}

			function getA() {
				return foo.a
			}
		`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		v, err := vm.RunFunc("getA")
		if err != nil {
			log.Fatal(err)
		}

		if v.ToInt() != 3 {
			log.Fatal(v)
		}
	}
}

func BenchmarkObjectField(b *testing.B) {
	vm := initVM(b, `
			let foo = { a: 3 }

			function getA() {
				return foo.a
			}
		`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		v, err := vm.RunFunc("getA")
		if err != nil {
			log.Fatal(err)
		}

		if v.ToInt() != 3 {
			log.Fatal(v)
		}
	}
}

func BenchmarkClassMethod(b *testing.B) {
	vm := initVM(b, `
			let foo = new Foo()

			class Foo {
				bar() {
					return 3
				}
			}

			function getA() {
				return foo.bar()
			}
		`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		v, err := vm.RunFunc("getA")
		if err != nil {
			log.Fatal(err)
		}

		if v.ToInt() != 3 {
			log.Fatal(v)
		}
	}
}

func BenchmarkObjectMethod(b *testing.B) {
	vm := initVM(b, `
			let foo = { bar: () => 3 }

			function getA() {
				return foo.bar()
			}
		`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		v, err := vm.RunFunc("getA")
		if err != nil {
			log.Fatal(err)
		}

		if v.ToInt() != 3 {
			log.Fatal(v)
		}
	}
}
