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

A web server:
```typescript
let s = http.newServer()
s.address = ":8080"
s.handler = (w, r) => w.write("Hello world")
s.start() 
```

Array built in functions:
```typescript
let items = [1, 2, 3, 4, 5]
let v = items.where(t => t > 2).select(t => t + 3).sum()
console.log(v)
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