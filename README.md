# Dune 

## Install 
```
go get github.com/scorredoira/dune/cmd/dune
```


## Write some code

Open hello.ts in your editor and type a program
```typescript
fmt.println("Hello, World!")
```

Run your code
```
$ dune hello.ts
Hello, World!
```


The syntax is a subset of Typescript. This allows to get type checking, autocomplete and refactoring support from editors like VSCode.

Generate a project and type definitions 

```
$ dune -init
```

## Examples

Array built in functions:
```typescript
let items = [1, 2, 3, 4, 5]
let v = items.where(t => t > 2).select(t => t + 3).sum()
console.log(v)
```

A web server:
```typescript
let s = http.newServer()
s.address = ":8080"
s.handler = (w, r) => w.write("Hello world")
s.start() 
```

With autocert:
```typescript
let tlsconf = tls.newConfig()
tlsconf.certManager = autocert.newCertManager("certs", ["example.com"])

let s = http.newServer()
s.tlsConfig = tlsconf
s.handler = (w, r) => w.write("Hello world")
s.start()
```

Checkout more [examples](https://github.com/scorredoira/dune-examples).


## REPL
```
$ dune
DUNE 0.93 Build
commands: :paste, list, asm, quit

> 1 + 2
3
```



## Embedding

```Go
package main

import (
	"fmt"

	"github.com/scorredoira/dune"
)

func main() {
	v, err := dune.RunStr("return 3 * 2")
	fmt.Println(v, err)
}
```


## Cross compile from linux.

CGO is necessary for sqlite support.

To install a Windows C compiler in Ubuntu

```
sudo apt-get install mingw-w64
```

```
GOOS=darwin GOARCH=amd64 go build -o dune-mac cmd/dune/main.go

GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -o dune.exe cmd/dune/main.go

GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i586-mingw32msvc-gcc go build -o dune.exe cmd/dune/main.go
```

