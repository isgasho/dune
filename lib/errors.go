package lib

import (
	"errors"
	"fmt"

	"github.com/scorredoira/dune"
)

var ErrReadOnly = errors.New("readonly property")
var ErrReadOnlyOrUndefined = errors.New("undefined or readonly property")
var ErrUndefined = errors.New("undefined")
var ErrInvalidType = errors.New("invalid value type")
var ErrFileNotFound = errors.New("file not found")
var ErrUnauthorized = errors.New("unauthorized")
var ErrNoFileSystem = errors.New("there is no filesystem")

func init() {
	dune.RegisterLib(Errors, `

	declare namespace errors {
		export function newError(msg: string): Error
		export function wrap(msg: string, inner: Error): Error
		export function public(msg: string, inner?: Error | string): Error
	
		export function is(err: Error, msg: string): Error
		export function rethrow(err: Error): void

		export interface Error {
			public: boolean
			message: string
			pc: number
			stackTrace: string
			toString(): string
			is(error: string): boolean
		}
	}
`)
}

var Errors = []dune.NativeFunction{
	{
		Name:      "errors.is",
		Arguments: 2,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if err := ValidateArgs(args, dune.Object, dune.String); err != nil {
				return dune.NullValue, err
			}
			e, ok := args[0].ToObjectOrNil().(dune.Error)
			if !ok {
				return dune.FalseValue, nil
			}
			return dune.NewBool(e.Is(args[1].ToString())), nil
		},
	},
	{
		Name:      "errors.rethrow",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if err := ValidateArgs(args, dune.Object); err != nil {
				return dune.NullValue, err
			}

			e, ok := args[0].ToObjectOrNil().(dune.Error)
			if !ok {
				return dune.NullValue, fmt.Errorf("Expected error, got %s", args[0].ToString())
			}

			e.IsRethrow = true

			return dune.NullValue, e
		},
	},
	{
		Name:      "errors.newError",
		Arguments: 1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if err := ValidateArgs(args, dune.String); err != nil {
				return dune.NullValue, err
			}
			return dune.NewObject(vm.NewError(args[0].ToString())), nil
		},
	},
	{
		Name:      "errors.wrap",
		Arguments: -1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			return wrap(false, args, vm)
		},
	},
	{
		Name:      "errors.public",
		Arguments: -1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			return wrap(true, args, vm)
		},
	},
}

func wrap(public bool, args []dune.Value, vm *dune.VM) (dune.Value, error) {
	ln := len(args)
	if ln < 1 || ln > 2 {
		return dune.NullValue, fmt.Errorf("expected 1 or 2 parameters, got %d", ln)
	}

	v := args[0]
	if v.Type != dune.String {
		return dune.NullValue, fmt.Errorf("expected parameter 1 to be a string, got %s", v.Type)
	}

	e := vm.NewPublicError(v.ToString())

	if ln > 1 {
		innerObj := args[1]
		switch innerObj.Type {

		case dune.Null, dune.Undefined:

		case dune.String:
			innerEx := vm.NewError(innerObj.ToString())
			e.Wrap(innerEx)

		case dune.Object:
			if innerObj.Type != dune.Object {
				return dune.NullValue, fmt.Errorf("expected parameter 2 to be a Exception, got %s", innerObj.Type)
			}
			innerEx, ok := innerObj.ToObject().(dune.Error)
			if !ok {
				return dune.NullValue, fmt.Errorf("expected parameter 2 to be a Exception, got %s", innerEx.Type())
			}
			e.Wrap(innerEx)

		default:
			return dune.NullValue, fmt.Errorf("expected parameter 2 to be a Exception, got %s", innerObj.Type)
		}
	}

	return dune.NewObject(e), nil
}
