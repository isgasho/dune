package lib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/scorredoira/dune"
	"github.com/scorredoira/dune/binary"
)

func init() {
	dune.RegisterLib(Runtime, `

declare function panic(message: string): void
declare function defer(f: () => void): void;

declare namespace runtime {	
    export const version: string
}

declare namespace runtime {
	export const ErrFunctionNotExist: string

    export interface Finalizable { 
        close(): void
	}
	
	// export function call(module: string, func: string, ...args: any[]): any

	export function typeDefs(): string

    export function setFileSystem(fs: io.FileSystem): void

    export function setFinalizer(v: runtime.Finalizable): void
    export function newFinalizable(f: () => void): Finalizable

    export function panic(message: string): void

    export type OSName = "linux" | "windows" | "darwin"

	export const context: any
	
    /**
     * Returns the operating system
     */
    export const OS: OSName

    /**
     * Returns the path of the executable.
     */
    export const executable: string

    /**
     * Returns the path of the native runtime executable.
     */
    export const nativeExecutable: string

    export const vm: VirtualMachine


    export function runFunc(func: string, ...args: any[]): any

    export const hasResources: boolean
    export function resources(name: string): string[]
    export function resource(name: string): byte[]

    export function getStackTrace(): string
    export function newVM(p: Program, globals?: any[]): VirtualMachine

    export interface Program {
		readonly constants: any[]
        functions(): FunctionInfo[]
        functionInfo(name: string): FunctionInfo
        resources(): string[]
        resource(key: string): byte[]
        setResource(key: string, value: byte[]): void

		directives(): string[]
		directive(): string
		hasDirective(name: string): boolean

		permissions(): string[]
		hasPermission(name: string): boolean
		addPermission(name: string): void

        /**
         * Strip sources, not exported functions name and other info.
         */
        strip(): void
        toString(): string
        write(w: io.Writer): void
	}
	
    export interface FunctionInfo {
        name: string
        index: number
		arguments: number
		optionalArguments: number
        exported: boolean
		func: Function
		directives(): string[]
		directive(): string
		hasDirective(name: string): boolean
        toString(): string
    }

    export interface VirtualMachine {
		maxAllocations: number
		maxFrames: number
		maxSteps: number
		fileSystem: io.FileSystem
		localizer: locale.Localizer
		readonly steps: number
		readonly allocations: number
		readonly program: Program
		context: any
		error: errors.Error
		initialize(): any[]
		run(...args: any[]): any
		runFunc(name: string, ...args: any[]): any
		runFunc(index: number, ...args: any[]): any
		runStaticFunc(name: string, ...args: any[]): any
		runStaticFunc(index: number, ...args: any[]): any
		getValue(name: string): any
		getGlobals(): any[]
		getStackTrace(): string
		clone(): VirtualMachine
		resetSteps(): void
	}
}


`)
}

var Runtime = []dune.NativeFunction{
	{
		Name: "->runtime.ErrFunctionNotExist",
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			return dune.NewString(dune.ErrFunctionNotExist.Error()), nil
		},
	},
	{
		Name:      "panic",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if !vm.HasPermission("trusted") {
				return dune.NullValue, ErrUnauthorized
			}

			panic(args[0].ToString())
		},
	},
	{
		Name: "->runtime.version",
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			return dune.NewString(dune.VERSION), nil
		},
	},
	{
		Name: "runtime.typeDefs",
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if err := ValidateArgs(args); err != nil {
				return dune.NullValue, err
			}
			s := dune.TypeDefs()
			return dune.NewString(s), nil
		},
	},
	{
		Name:      "runtime.runFunc",
		Arguments: -1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if !vm.HasPermission("trusted") {
				return dune.NullValue, ErrUnauthorized
			}

			if len(args) == 0 {
				return dune.NullValue, fmt.Errorf("expected at least the function name")
			}

			a := args[0]
			if a.Type != dune.String {
				return dune.NullValue, fmt.Errorf("function name must be a string, got %v", a.Type)
			}

			funcName := a.ToString()

			return vm.RunFunc(funcName, args[1:]...)
		},
	},
	{
		Name: "->runtime.resources",
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			res := vm.Program.Resources
			if res == nil {
				return dune.NewArray(0), nil
			}

			a := make([]dune.Value, len(res))

			i := 0
			for k := range res {
				a[i] = dune.NewString(k)
				i++
			}

			return dune.NewArrayValues(a), nil
		},
	},
	{
		Name:      "runtime.resource",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if err := ValidateArgs(args, dune.String); err != nil {
				return dune.NullValue, err
			}

			name := args[0].ToString()

			res := vm.Program.Resources
			if res == nil {
				return dune.NullValue, nil
			}

			v, ok := res[name]
			if !ok {
				return dune.NullValue, nil
			}

			return dune.NewBytes(v), nil
		},
	},
	{
		Name: "->runtime.hasResources",
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			res := vm.Program.Resources
			if len(res) == 0 {
				return dune.FalseValue, nil
			}

			return dune.TrueValue, nil
		},
	},
	{
		Name: "->runtime.context",
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if !vm.HasPermission("trusted") {
				return dune.NullValue, ErrUnauthorized
			}
			return vm.Context, nil
		},
	},
	{
		Name:      "runtime.setFileSystem",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if !vm.HasPermission("trusted") {
				return dune.NullValue, ErrUnauthorized
			}

			if err := ValidateArgs(args, dune.Object); err != nil {
				return dune.NullValue, err
			}

			fs, ok := args[0].ToObject().(*FileSystemObj)
			if !ok {
				return dune.NullValue, fmt.Errorf("expected a fileSystem, got %s", args[0].TypeName())
			}
			vm.FileSystem = fs.FS
			return dune.NullValue, nil
		},
	},
	{
		Name:      "runtime.newFinalizable",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			v := args[0]

			fin, err := newFinalizable(v, vm)
			if err != nil {
				return dune.NullValue, err
			}

			return dune.NewObject(fin), nil
		},
	},
	{
		Name:      "defer",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			v := args[0]

			fin, err := newFinalizable(v, vm)
			if err != nil {
				return dune.NullValue, err
			}

			vm.SetFinalizer(fin)
			return dune.NullValue, nil
		},
	},
	{
		Name:      "runtime.setFinalizer",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			v := args[0]

			if v.Type != dune.Object {
				return dune.NullValue, fmt.Errorf("the value is not a finalizer")
			}

			fin, ok := v.ToObject().(dune.Finalizable)
			if !ok {
				return dune.NullValue, fmt.Errorf("the value is not a finalizer")
			}
			vm.SetFinalizer(fin)
			return dune.NullValue, nil
		},
	},
	{
		Name: "->runtime.OS",
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if !vm.HasPermission("trusted") {
				return dune.NullValue, ErrUnauthorized
			}
			return dune.NewString(runtime.GOOS), nil
		},
	},
	{
		Name: "->runtime.nativeExecutable",
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if !vm.HasPermission("trusted") {
				return dune.NullValue, ErrUnauthorized
			}
			ex, err := os.Executable()
			if err != nil {
				return dune.NullValue, err
			}
			return dune.NewString(ex), nil
		},
	},
	{
		Name:      "runtime.newVM",
		Arguments: -1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			l := len(args)
			if l == 0 || l > 2 {
				return dune.NullValue, fmt.Errorf("expected 1 or 2 params, got %d", l)
			}

			if args[0].Type != dune.Object {
				return dune.NullValue, fmt.Errorf("argument 1 must be a program, got %s", args[0].TypeName())
			}
			p, ok := args[0].ToObject().(*program)
			if !ok {
				return dune.NullValue, fmt.Errorf("argument 1 must be a program, got %s", args[0].TypeName())
			}

			var m *dune.VM

			if l == 1 {
				m = dune.NewVM(p.prog)
			} else {
				switch args[1].Type {
				case dune.Undefined, dune.Null:
					m = dune.NewVM(p.prog)
				case dune.Array:
					m = dune.NewInitializedVM(p.prog, args[1].ToArray())
				default:
					return dune.NullValue, fmt.Errorf("argument 2 must be an array, got %s", args[1].TypeName())
				}
			}

			m.MaxAllocations = vm.MaxAllocations
			m.MaxFrames = vm.MaxFrames
			m.MaxSteps = vm.MaxSteps

			if err := m.AddSteps(vm.Steps()); err != nil {
				return dune.NullValue, err
			}

			return dune.NewObject(&libVM{m}), nil
		},
	},
	{
		Name: "->runtime.vm",
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if !vm.HasPermission("trusted") {
				return dune.NullValue, ErrUnauthorized
			}
			return dune.NewObject(&libVM{vm}), nil
		},
	},
	{
		Name: "runtime.resetSteps",
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if !vm.HasPermission("trusted") {
				return dune.NullValue, ErrUnauthorized
			}
			vm.ResetSteps()
			return dune.NullValue, nil
		},
	},
	{
		Name:      "runtime.getStackTrace",
		Arguments: 0,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			s := strings.Join(vm.Stacktrace(), "\n")
			return dune.NewString(s), nil
		},
	},
}

func newFinalizable(v dune.Value, vm *dune.VM) (finalizable, error) {
	switch v.Type {
	case dune.Func:

	case dune.NativeFunc:

	case dune.Object:
		switch v.ToObject().(type) {
		case *dune.Closure, dune.NativeMethod:
		default:
			return finalizable{}, fmt.Errorf("expected a function, got: %s", v.TypeName())
		}

	default:
		return finalizable{}, fmt.Errorf("expected a function, got %v", v.TypeName())
	}

	f := finalizable{v: v, vm: vm}
	return f, nil
}

type finalizable struct {
	v  dune.Value
	vm *dune.VM
}

func (finalizable) Type() string {
	return "[Finalizable]"
}

func (f finalizable) Close() error {
	v := f.v
	vm := f.vm

	var lastErr = vm.Error
	if lastErr != nil {
		defer func() {
			vm.Error = lastErr
		}()
	}

	vm.Error = nil
	switch v.Type {

	case dune.NativeFunc:
		i := v.ToNativeFunction()
		f := dune.NativeFuncFromIndex(i)
		if f.Arguments != 0 {
			return fmt.Errorf("function '%s' expects %d parameters", f.Name, f.Arguments)
		}
		_, err := f.Function(dune.NullValue, nil, vm)
		return err

	case dune.Func:
		i := v.ToFunction()
		if _, err := vm.RunFuncIndex(i); err != nil {
			return err
		}

	case dune.Object:
		switch t := v.ToObject().(type) {
		case *dune.Closure:
			if _, err := vm.RunClosure(t); err != nil {
				return err
			}
		case dune.NativeMethod:
			if _, err := t(nil, vm); err != nil {
				return err
			}
		default:
			panic(fmt.Sprintf("expected a function, got: %s", v.TypeName()))
		}

	default:
		panic("should be a function or a closure")

	}

	return nil
}

func (f finalizable) GetMethod(name string) dune.NativeMethod {
	switch name {
	case "close":
		return f.close
	}
	return nil
}

func (f finalizable) close(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	return dune.NullValue, nil
}

type program struct {
	prog *dune.Program
}

func (p *program) Type() string {
	return "runtime.Program"
}

func (p *program) GetProperty(name string, vm *dune.VM) (dune.Value, error) {
	switch name {
	case "constants":
		return dune.NewArrayValues(p.prog.Constants), nil
	}

	return dune.UndefinedValue, nil
}

func (p *program) GetMethod(name string) dune.NativeMethod {
	switch name {
	case "functions":
		return p.functions
	case "functionInfo":
		return p.functionInfo
	case "toString":
		return p.toString
	case "toBytes":
		return p.toBytes
	case "resources":
		return p.resources
	case "setResource":
		return p.setResource
	case "resource":
		return p.resource
	case "strip":
		return p.strip
	case "write":
		return p.write
	case "directives":
		return p.directives
	case "directive":
		return p.directive
	case "permissions":
		return p.permissions
	case "hasDirective":
		return p.hasDirective
	case "hasPermission":
		return p.hasPermission
	case "addPermission":
		return p.addPermission
	}
	return nil
}

func (p *program) hasPermission(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args, dune.String); err != nil {
		return dune.NullValue, err
	}

	name := args[0].ToString()

	v := p.prog.HasPermission(name)

	return dune.NewBool(v), nil
}

func (p *program) addPermission(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if !vm.HasPermission("trusted") {
		return dune.NullValue, ErrUnauthorized
	}

	if err := ValidateArgs(args, dune.String); err != nil {
		return dune.NullValue, err
	}

	name := args[0].ToString()
	p.prog.AddPermission(name)

	return dune.NullValue, nil
}

func (p *program) hasDirective(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args, dune.String); err != nil {
		return dune.NullValue, err
	}

	name := args[0].ToString()

	var found bool

	for _, v := range p.prog.Directives {
		if v == name {
			found = true
			break
		}
	}
	return dune.NewBool(found), nil
}

func (p *program) directives(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args); err != nil {
		return dune.NullValue, err
	}

	result := make([]dune.Value, len(p.prog.Directives))
	for i, item := range p.prog.Directives {
		result[i] = dune.NewString(item)
	}
	return dune.NewArrayValues(result), nil
}

func (p *program) directive(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args, dune.String); err != nil {
		return dune.NullValue, err
	}

	name := args[0].ToString() + " "

	for _, directive := range p.prog.Directives {
		if strings.HasPrefix(directive, name) {
			return dune.NewString(strings.TrimPrefix(directive, name)), nil
		}
	}
	return dune.NullValue, nil
}

func (p *program) permissions(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args); err != nil {
		return dune.NullValue, err
	}

	permissions := p.prog.Permissions()

	result := make([]dune.Value, len(permissions))
	for i, item := range permissions {
		result[i] = dune.NewString(item)
	}
	return dune.NewArrayValues(result), nil
}

func (p *program) write(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if !vm.HasPermission("trusted") {
		return dune.NullValue, ErrUnauthorized
	}

	if err := ValidateArgs(args, dune.Object); err != nil {
		return dune.NullValue, err
	}

	w, ok := args[0].ToObjectOrNil().(io.Writer)
	if !ok {
		return dune.NullValue, fmt.Errorf("exepected a Writer, got %s", args[0].TypeName())
	}

	if err := binary.Write(w, p.prog); err != nil {
		return dune.NullValue, err
	}

	return dune.NullValue, nil
}

func (p *program) strip(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if !vm.HasPermission("trusted") {
		return dune.NullValue, ErrUnauthorized
	}

	if err := ValidateArgs(args); err != nil {
		return dune.NullValue, err
	}

	p.prog.Strip()

	return dune.NullValue, nil
}

func (p *program) setResource(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if !vm.HasPermission("trusted") {
		return dune.NullValue, ErrUnauthorized
	}

	if err := ValidateArgs(args, dune.String, dune.Bytes); err != nil {
		return dune.NullValue, err
	}

	if p.prog.Resources == nil {
		p.prog.Resources = make(map[string][]byte)
	}

	p.prog.Resources[args[0].ToString()] = args[1].ToBytes()
	return dune.NullValue, nil
}

func (p *program) resources(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args); err != nil {
		return dune.NullValue, err
	}

	res := p.prog.Resources

	if res == nil {
		return dune.NewArray(0), nil
	}

	a := make([]dune.Value, len(res))

	i := 0
	for k := range res {
		a[i] = dune.NewString(k)
		i++
	}

	return dune.NewArrayValues(a), nil
}

func (p *program) resource(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args, dune.String); err != nil {
		return dune.NullValue, err
	}

	name := args[0].ToString()

	if p.prog.Resources == nil {
		return dune.NullValue, nil
	}

	v, ok := p.prog.Resources[name]
	if !ok {
		return dune.NullValue, nil
	}

	return dune.NewBytes(v), nil
}

func (p *program) functions(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if len(args) != 0 {
		return dune.NullValue, fmt.Errorf("expected no args")
	}

	var funcs []dune.Value
	for _, f := range p.prog.Functions {
		fi := functionInfo{f, *p}
		funcs = append(funcs, dune.NewObject(fi))
	}
	return dune.NewArrayValues(funcs), nil
}

func (p *program) functionInfo(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args, dune.String); err != nil {
		return dune.NullValue, err
	}

	name := args[0].ToString()

	f, ok := p.prog.Function(name)
	if !ok {
		return dune.NullValue, nil
	}

	return dune.NewObject(functionInfo{f, *p}), nil
}

func (p *program) toBytes(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args, dune.Object); err != nil {
		return dune.NullValue, err
	}

	w, ok := args[0].ToObject().(io.Writer)
	if !ok {
		return dune.NullValue, fmt.Errorf("expected parameter 1 to be io.Writer, got %T", args[0].ToObject())
	}

	err := binary.Write(w, p.prog)
	return dune.NullValue, err
}

func (p *program) toString(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	var b bytes.Buffer
	dune.Fprint(&b, p.prog)
	return dune.NewString(b.String()), nil
}

type functionInfo struct {
	fn *dune.Function
	p  program
}

func (functionInfo) Type() string {
	return "runtime.FunctionInfo"
}

func (f functionInfo) GetMethod(name string) dune.NativeMethod {
	switch name {
	case "directives":
		return f.directives
	case "directive":
		return f.directive
	case "hasDirective":
		return f.hasDirective
	case "toString":
		return f.toString
	}
	return nil
}

func (f functionInfo) hasDirective(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args, dune.String); err != nil {
		return dune.NullValue, err
	}

	name := args[0].ToString()

	var found bool

	for _, v := range f.fn.Directives {
		if v == name {
			found = true
			break
		}
	}
	return dune.NewBool(found), nil
}

func (f functionInfo) directives(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args); err != nil {
		return dune.NullValue, err
	}

	result := make([]dune.Value, len(f.fn.Directives))
	for i, item := range f.fn.Directives {
		result[i] = dune.NewString(item)
	}
	return dune.NewArrayValues(result), nil
}

func (f functionInfo) directive(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if err := ValidateArgs(args, dune.String); err != nil {
		return dune.NullValue, err
	}

	name := args[0].ToString() + " "

	for _, directive := range f.fn.Directives {
		if strings.HasPrefix(directive, name) {
			return dune.NewString(strings.TrimPrefix(directive, name)), nil
		}
	}
	return dune.NullValue, nil
}

func (f functionInfo) toString(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	var b bytes.Buffer
	dune.FprintFunction(&b, f.fn, f.p.prog)
	return dune.NewString(b.String()), nil
}

func (f functionInfo) GetProperty(name string, vm *dune.VM) (dune.Value, error) {
	switch name {
	case "name":
		return dune.NewString(f.fn.Name), nil
	case "arguments":
		return dune.NewInt(f.fn.Arguments), nil
	case "optionalArguments":
		return dune.NewInt(f.fn.OptionalArguments), nil
	case "index":
		return dune.NewInt(f.fn.Index), nil
	case "exported":
		return dune.NewBool(f.fn.Exported), nil
	case "func":
		return dune.NewFunction(f.fn.Index), nil
	}
	return dune.UndefinedValue, nil
}

type libVM struct {
	vm *dune.VM
}

func (m *libVM) Type() string {
	return "runtime.VirtualMachine"
}

func (m *libVM) GetProperty(name string, vm *dune.VM) (dune.Value, error) {
	switch name {
	case "context":
		return vm.Context, nil
	case "error":
		e := m.vm.Error
		if e != nil {
			return dune.NewObject(e), nil
		}
		return dune.NullValue, nil
	case "program":
		if !vm.HasPermission("trusted") {
			return dune.NullValue, ErrUnauthorized
		}
		return dune.NewObject(&program{prog: m.vm.Program}), nil
	case "fileSystem":
		return dune.NewObject(NewFileSystem(m.vm.FileSystem)), nil
	case "language":
		return dune.NewString(m.vm.Language), nil
	case "localizer":
		return dune.NewObject(m.vm.Localizer), nil
	case "maxAllocations":
		return dune.NewInt64(m.vm.MaxAllocations), nil
	case "maxFrames":
		return dune.NewInt(m.vm.MaxFrames), nil
	case "maxSteps":
		return dune.NewInt64(m.vm.MaxSteps), nil
	case "steps":
		return dune.NewInt64(m.vm.Steps()), nil
	case "allocations":
		return dune.NewInt64(m.vm.Allocations()), nil
	}
	return dune.UndefinedValue, nil
}

func (m *libVM) SetProperty(name string, v dune.Value, vm *dune.VM) error {
	if !vm.HasPermission("trusted") {
		return ErrUnauthorized
	}

	switch name {
	case "error":
		switch v.Type {
		case dune.Null:
			m.vm.Error = nil
			return nil

		case dune.Object:
			e, ok := v.ToObject().(error)
			if !ok {
				return ErrInvalidType
			}
			m.vm.Error = e
			return nil

		default:
			return ErrInvalidType
		}

	case "context":
		m.vm.Context = v
		return nil

	case "fileSystem":
		fs, ok := v.ToObjectOrNil().(*FileSystemObj)
		if !ok {
			return ErrInvalidType
		}
		m.vm.FileSystem = fs.FS
		return nil

	case "language":
		if v.Type != dune.String {
			return ErrInvalidType
		}
		m.vm.Language = v.ToString()
		return nil

	case "localizer":
		loc, ok := v.ToObjectOrNil().(dune.Localizer)
		if !ok {
			return ErrInvalidType
		}
		m.vm.Localizer = loc
		return nil

	case "maxAllocations":
		if v.Type != dune.Int {
			return ErrInvalidType
		}
		m.vm.MaxAllocations = v.ToInt()
		return nil

	case "maxFrames":
		if v.Type != dune.Int {
			return ErrInvalidType
		}
		m.vm.MaxFrames = int(v.ToInt())
		return nil

	case "maxSteps":
		if v.Type != dune.Int {
			return ErrInvalidType
		}
		m.vm.MaxSteps = v.ToInt()
		return nil
	}

	return ErrReadOnlyOrUndefined
}

func (m *libVM) GetMethod(name string) dune.NativeMethod {
	switch name {
	case "initialize":
		return m.initialize
	case "run":
		return m.run
	case "runFunc":
		return m.runFunc
	case "runStaticFunc":
		return m.runStaticFunc
	case "clone":
		return m.clone
	case "getValue":
		return m.getValue
	case "getGlobals":
		return m.getGlobals
	case "getStackTrace":
		return m.getStackTrace
	case "resetSteps":
		return m.resetSteps
	}
	return nil
}

func (m *libVM) clone(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if !vm.HasPermission("trusted") {
		return dune.NullValue, ErrUnauthorized
	}

	if err := ValidateArgs(args); err != nil {
		return dune.NullValue, err
	}

	c := m.vm.Clone(m.vm.Program, m.vm.Globals())
	return dune.NewObject(&libVM{c}), nil
}

func (m *libVM) resetSteps(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if !vm.HasPermission("trusted") {
		return dune.NullValue, ErrUnauthorized
	}

	m.vm.ResetSteps()

	return dune.NullValue, nil
}

func (m *libVM) getStackTrace(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if !vm.HasPermission("trusted") {
		return dune.NullValue, ErrUnauthorized
	}

	s := strings.Join(m.vm.Stacktrace(), "\n")
	return dune.NewString(s), nil
}

func (m *libVM) getGlobals(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if !vm.HasPermission("trusted") {
		return dune.NullValue, ErrUnauthorized
	}

	return dune.NewArrayValues(m.vm.Globals()), nil
}

func (m *libVM) initialize(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if !vm.HasPermission("trusted") {
		return dune.NullValue, ErrUnauthorized
	}

	if err := m.vm.Initialize(); err != nil {
		// return the error with the stacktrace included in the message
		// because the caller in the program will have it's own stacktrace.
		return dune.NullValue, errors.New(err.Error())
	}

	return dune.NewArrayValues(m.vm.Globals()), nil
}

func (m *libVM) run(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if !vm.HasPermission("trusted") {
		return dune.NullValue, ErrUnauthorized
	}

	v, err := m.vm.Run(args...)
	if err != nil {
		// return the error with the stacktrace included in the message
		// because the caller in the program will have it's own stacktrace.
		return dune.NullValue, errors.New(err.Error())
	}
	return v, nil
}

func (m *libVM) runStaticFunc(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	return m.runFunction(args, vm, false)
}

func (m *libVM) runFunc(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	return m.runFunction(args, vm, true)
}

func (m *libVM) runFunction(args []dune.Value, vm *dune.VM, initialize bool) (dune.Value, error) {
	if !vm.HasPermission("trusted") {
		return dune.NullValue, ErrUnauthorized
	}

	l := len(args)
	if l == 0 {
		return dune.NullValue, fmt.Errorf("expected at least 1 parameter, got %d", l)
	}

	var index int

	switch args[0].Type {
	case dune.String:
		name := args[0].ToString()
		f, ok := m.vm.Program.Function(name)
		if !ok {
			return dune.NullValue, fmt.Errorf("%s: %w", name, dune.ErrFunctionNotExist)
		}
		index = f.Index
	case dune.Int:
		index = int(args[0].ToInt())
		if len(m.vm.Program.Functions) <= index {
			return dune.NullValue, fmt.Errorf("%d: %w", index, dune.ErrFunctionNotExist)
		}
	default:
		return dune.NullValue, fmt.Errorf("argument 1 must be a string (function name), got %s", args[0].TypeName())
	}

	if initialize && !m.vm.Initialized() {
		if err := m.vm.Initialize(); err != nil {
			return dune.NullValue, m.vm.WrapError(err)
		}
		if err := vm.AddSteps(m.vm.Steps()); err != nil {
			return dune.NullValue, m.vm.WrapError(err)
		}
	}

	v, err := m.vm.RunFuncIndex(index, args[1:]...)
	if err != nil {
		// return the error with the stacktrace included in the message
		// because the caller in the program will have it's own stacktrace.
		return dune.NullValue, m.vm.WrapError(err)
	}
	if err := vm.AddSteps(m.vm.Steps()); err != nil {
		return dune.NullValue, m.vm.WrapError(err)
	}
	return v, nil
}

func (m *libVM) getValue(args []dune.Value, vm *dune.VM) (dune.Value, error) {
	if !vm.HasPermission("trusted") {
		return dune.NullValue, ErrUnauthorized
	}

	l := len(args)
	if l != 1 {
		return dune.NullValue, fmt.Errorf("expected 1 parameter, got %d", l)
	}

	if args[0].Type != dune.String {
		return dune.NullValue, fmt.Errorf("argument 1 must be a string (var name), got %s", args[0].TypeName())
	}

	name := args[0].ToString()
	v, _ := m.vm.RegisterValue(name)
	return v, nil
}
