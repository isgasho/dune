package lib

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/scorredoira/dune"
)

func init() {
	dune.RegisterLib(Assert, `

declare namespace assert {
    export function contains(value: string, search: string): void
    export function equal(a: any, b: any): void
    export function isNull(a: any): void
	export function isNotNull(a: any): void
	export function exception(msg: string, func: Function): void

	export function int(a: any, msg: string): number
	export function float(a: any, msg: string): number
	export function string(a: any, msg: string): string
	export function bool(a: any, msg: string): boolean	
	export function object(a: any, msg: string): any	
}

`)
}

var Assert = []dune.NativeFunction{
	{
		Name:      "assert.contains",
		Arguments: 2,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if err := ValidateArgs(args, dune.String, dune.String); err != nil {
				return dune.NullValue, err
			}

			a := args[0].ToString()
			b := args[1].ToString()

			if !strings.Contains(a, b) {
				return dune.NullValue, fmt.Errorf("'%s' not contained in '%s'", b, a)
			}
			return dune.NullValue, nil
		},
	},

	{
		Name:      "assert.equal",
		Arguments: 2,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			a := args[0]
			b := args[1]

			if !areEqual(a, b) {
				return dune.NullValue, fmt.Errorf("values are different: %v, %v", serializeOrErr(a), serializeOrErr(b))
			}
			return dune.NullValue, nil
		},
	},
	{
		Name:      "assert.isNull",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			a := args[0]

			switch a.Type {
			case dune.Null, dune.Undefined:
			default:
				return dune.NullValue, fmt.Errorf("expected null, got %v", a)
			}

			return dune.NullValue, nil
		},
	},
	{
		Name:      "assert.isNotNull",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			a := args[0]

			switch a.Type {
			case dune.Null, dune.Undefined:
				return dune.NullValue, fmt.Errorf("%v is null", a)
			default:
			}

			return dune.NullValue, nil
		},
	},
	{
		Name:      "assert.exception",
		Arguments: 2,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			a := args[0]
			if a.Type != dune.String {
				return dune.NullValue, fmt.Errorf("expected argument 1 to be a string, got %s", a.TypeName())
			}

			expected := a.ToString()

			v := args[1]
			err := runFuncOrClosure(vm, v)
			if err == nil {
				return dune.NullValue, fmt.Errorf("expected exception: %s", expected)
			}

			if expected != "" && !strings.Contains(err.Error(), expected) {
				return dune.NullValue, fmt.Errorf("invalid exception, does not contain '%s': %s", expected, err.Error())
			}

			// clear the error
			vm.Error = nil

			return dune.NullValue, nil
		},
	},
	{
		Name:      "assert.int",
		Arguments: 2,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if args[1].Type != dune.String {
				return dune.NullValue, fmt.Errorf("expected argument 2 to be a string, got %s", args[1].TypeName())
			}

			a := args[0]
			msg := args[1].ToString()

			var v int64
			var err error

			switch a.Type {
			case dune.Int:
				v = a.ToInt()
			case dune.String:
				v, err = strconv.ParseInt(a.ToString(), 0, 64)
				if err != nil {
					return dune.NullValue, fmt.Errorf(msg, showAssertMessage("%v is not int", a.TypeName()))
				}
			default:
				return dune.NullValue, fmt.Errorf(msg)
			}

			return dune.NewInt64(v), nil
		},
	},
	{
		Name:      "assert.float",
		Arguments: 2,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if args[1].Type != dune.String {
				return dune.NullValue, fmt.Errorf("expected argument 2 to be a string, got %s", args[1].TypeName())
			}

			a := args[0]
			msg := args[1].ToString()

			var v int64
			var err error

			switch a.Type {
			case dune.Int:
				v = a.ToInt()
			case dune.String:
				v, err = strconv.ParseInt(a.ToString(), 0, 64)
				if err != nil {
					return dune.NullValue, fmt.Errorf(msg, showAssertMessage("%v is not float", a.TypeName()))
				}
			default:
				return dune.NullValue, fmt.Errorf(msg, showAssertMessage("%v is not float", a.TypeName()))
			}

			return dune.NewInt64(v), nil
		},
	},
	{
		Name:      "assert.string",
		Arguments: 2,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if args[1].Type != dune.String {
				return dune.NullValue, fmt.Errorf("expected argument 2 to be a string, got %s", args[1].TypeName())
			}

			a := args[0]
			msg := args[1].ToString()

			var v string

			switch a.Type {
			case dune.Int, dune.Float, dune.Bool, dune.String:
				v = a.ToString()
			default:
				return dune.NullValue, fmt.Errorf(msg, showAssertMessage("%v is not a string", a.TypeName()))
			}

			return dune.NewString(v), nil
		},
	},
	{
		Name:      "assert.bool",
		Arguments: 2,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if args[1].Type != dune.String {
				return dune.NullValue, fmt.Errorf("expected argument 2 to be a string, got %s", args[1].TypeName())
			}

			a := args[0]
			msg := args[1].ToString()
			var v dune.Value

			switch a.Type {

			case dune.Bool:
				v = a

			case dune.Int:
				switch a.ToInt() {
				case 0:
					v = dune.FalseValue
				case 1:
					v = dune.TrueValue
				default:
					return dune.NullValue, fmt.Errorf(msg, showAssertMessage("%v is not bool", a.TypeName()))
				}

			case dune.String:
				s := a.ToString()
				s = strings.Trim(s, " ")
				switch s {
				case "true", "1":
					v = dune.TrueValue
				case "false", "0":
					v = dune.FalseValue
				default:
					return dune.NullValue, fmt.Errorf(msg, showAssertMessage("%v is not bool", a.TypeName()))
				}

			default:
				return dune.NullValue, fmt.Errorf(msg, showAssertMessage("%v is not bool", a.TypeName()))

			}

			return v, nil
		},
	},
	{
		Name:      "assert.object",
		Arguments: 2,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if args[1].Type != dune.String {
				return dune.NullValue, fmt.Errorf("expected argument 2 to be a string, got %s", args[1].TypeName())
			}

			a := args[0]
			msg := args[1].ToString()

			switch a.Type {
			case dune.Map:
			default:
				return dune.NullValue, fmt.Errorf(msg, showAssertMessage("%v is not an object", a.TypeName()))
			}

			return a, nil
		},
	},
}

func showAssertMessage(format string, args ...interface{}) string {
	if !strings.Contains(format, "%s") && !strings.Contains(format, "%v") {
		format += ": %s"
	}

	return fmt.Sprintf(format, args...)
}

func areEqual(a, b dune.Value) bool {
	if a.Equals(b) {
		return true
	}

	if a.Type == dune.Array && b.Type == dune.Array {
		aa := a.ToArrayObject().Array
		bb := b.ToArrayObject().Array
		if len(aa) != len(bb) {
			return false
		}
		for i, v := range aa {
			if !bb[i].Equals(v) {
				return false
			}
		}
		return true
	}

	return false
}

func serializeOrErr(v dune.Value) string {
	s, err := serialize(v)
	if err != nil {
		return err.Error()
	}
	return s
}
