package lib

import (
	"fmt"
	"io"

	"github.com/scorredoira/dune"
)

func init() {
	dune.RegisterLib(FileUtil, `

declare namespace fileutil { 
    export function copy(src: string, dst: string): byte[]
}

`)
}

var FileUtil = []dune.NativeFunction{
	{
		Name:      "fileutil.copy",
		Arguments: 2,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			if err := ValidateArgs(args, dune.String, dune.String); err != nil {
				return dune.NullValue, err
			}

			src := args[0].ToString()
			dst := args[1].ToString()

			fs := vm.FileSystem
			if fs == nil {
				return dune.NullValue, fmt.Errorf("no filesystem")
			}

			r, err := fs.Open(src)
			if err != nil {
				return dune.NullValue, err
			}

			w, err := fs.OpenForWrite(dst)
			if err != nil {
				return dune.NullValue, err
			}

			if _, err := io.Copy(w, r); err != nil {
				return dune.NullValue, err
			}

			return dune.NullValue, nil
		},
	},
}
