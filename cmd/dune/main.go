package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/scorredoira/dune"
	"github.com/scorredoira/dune/binary"
	"github.com/scorredoira/dune/filesystem"
	"github.com/scorredoira/dune/parser"

	_ "github.com/scorredoira/dune/lib"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func printVersion() {
	fmt.Printf("DUNE %s Build %s\n\n", dune.VERSION, dune.Build)
}

func main() {
	v := flag.Bool("v", false, "version")
	c := flag.Bool("c", false, "compile")
	o := flag.String("o", "", "output file")
	d := flag.Bool("d", false, "decompile")
	r := flag.Bool("r", false, "list resources")
	n := flag.Bool("n", false, "no optimizations")
	ini := flag.Bool("init", false, "generate native.d.ts and tsconfig.json")
	flag.Parse()

	if *v {
		printVersion()
		return
	}

	if *n {
		parser.Optimizations = false
	}

	args := flag.Args()
	aLen := len(args)

	if *d {
		p, err := loadProgram(args[0])
		if err != nil {
			fatal(err)
		}
		dune.Print(p)
		return
	}

	if *c {
		p, err := loadProgram(args[0])
		if err != nil {
			fatal(err)
		}

		out := *o
		if out == "" {
			n := filepath.Base(args[0])
			out = strings.TrimSuffix(n, filepath.Ext(n)) + ".bin"
		}
		if err := build(p, out); err != nil {
			fatal(err)
		}
		return
	}

	if *r {
		p, err := loadProgram(args[0])
		if err != nil {
			fatal(err)
		}
		for k, v := range p.Resources {
			fmt.Println(k, len(v))
		}
		return
	}

	if *ini {
		var path string
		if aLen == 1 {
			path = args[0]
		}
		generateDts(path)
		return
	}

	if aLen > 0 {
		if err := exec(args[0], args[1:]); err != nil {
			fatal(err)
		}
		return
	}

	if err := startREPL(); err != nil {
		fatal(err)
	}
}

func build(p *dune.Program, out string) error {
	f, err := os.OpenFile(out, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	if err := binary.Write(f, p); err != nil {
		return err
	}

	return nil
}

func exec(programPath string, args []string) error {
	p, err := loadProgram(programPath)
	if err != nil {
		return err
	}

	p.AddPermission("trusted")

	vm := dune.NewVM(p)
	vm.FileSystem = filesystem.OS

	ln := len(args)
	values := make([]dune.Value, ln)
	for i := 0; i < ln; i++ {
		values[i] = dune.NewValue(args[i])
	}

	_, err = vm.Run(values...)
	return err
}

func loadProgram(programPath string) (*dune.Program, error) {
	path, err := findPath(programPath)
	if err != nil {
		return nil, err
	}

	// by default source files have a typescript extension
	if strings.HasSuffix(path, ".ts") {
		return dune.Compile(filesystem.OS, path)
	}

	// first try to read as compiled
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening %s", path)
	}
	defer f.Close()

	p, err := binary.Read(f)
	if err != nil {
		if err == binary.ErrInvalidHeader {
			// if it is not a compiled program maybe is a source file with a different extension
			return dune.Compile(filesystem.OS, path)
		}
		return p, fmt.Errorf("error loading %s: %w", path, err)
	}

	return p, nil
}

func findPath(name string) (string, error) {
	if filepath.Ext(name) != "" {
		return name, nil
	}

	if path := tryPath(name); path != "" {
		return path, nil
	}

	// if it is just the name of the command search in envDirs
	if !strings.ContainsRune(name, os.PathSeparator) && filepath.Ext(name) == "" {
		path := os.Getenv("DUNE_DIRS")
		if path != "" {
			dirs := split(path, ":")
			for _, dir := range dirs {
				n := filepath.Join(dir, name)
				if path := tryPath(n); path != "" {
					return path, nil
				}
			}
		}
	}

	return "", fmt.Errorf("does not exist: %s", name)
}

func tryPath(name string) string {
	f, err := filesystem.OS.Stat(name)
	if err == nil && !f.IsDir() {
		return name
	}

	test := name + ".ts"
	if filesystem.Exists(filesystem.OS, test) {
		return test
	}

	test = name + ".bin"
	if filesystem.Exists(filesystem.OS, test) {
		return test
	}

	return ""
}

func generateDts(path string) {
	if path == "" {
		path = "."
	}

	filesystem.OS.WritePath(filepath.Join(path, "native.d.ts"), []byte(dune.TypeDefs()))

	writeIfNotExists(filesystem.OS, filepath.Join(path, "tsconfig.json"), []byte(`{
	"compilerOptions": {
		"noLib": true,
		"noEmit": true,
		"noImplicitAny": true,
		"baseUrl": "."
	}
}
`))
}

func writeIfNotExists(fs filesystem.FS, name string, data []byte) {
	if _, err := fs.Stat(name); err == nil {
		return
	}
	fs.WritePath(name, data)
}

func fatal(values ...interface{}) {
	fmt.Println(values...)
	os.Exit(1)
}

func split(s, sep string) []string {
	parts := strings.Split(s, sep)
	var result []string
	for _, p := range parts {
		if p != "" {
			// only append non empty values
			result = append(result, p)
		}
	}
	return result
}
