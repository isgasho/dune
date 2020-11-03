package benchmarks

// func BenchmarkNoInline(b *testing.B) {

// 	vm := initVM(b, `
// 			function main() {
// 				let v = 0
// 				for(let i = 0; i < 100; i++) {
// 					v += bar() + bar()
// 				}
// 				return v
// 			}

// 			function bar() {
// 				return 2 + 1
// 			}
// 		`)

// 	b.ResetTimer()
// 	b.ReportAllocs()

// 	for i := 0; i < b.N; i++ {
// 		v, err := vm.Run()
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		if v.ToInt() != 600 {
// 			log.Fatalf("got %v", v)
// 		}
// 	}
// }

// func BenchmarkInline(b *testing.B) {

// 	vm := initVM(b, `
// 			function main() {
// 				let v = 0
// 				for(let i = 0; i < 100; i++) {
// 					v += bar() + bar()
// 				}
// 				return v
// 			}

// 			function bar() {
// 				return 2 + 1
// 			}
// 		`)

// 	b.ResetTimer()
// 	b.ReportAllocs()

// 	for i := 0; i < b.N; i++ {
// 		v, err := vm.Run()
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		if v.ToInt() != 600 {
// 			log.Fatalf("got %v", v)
// 		}
// 	}
// }
