package lib

import (
	"encoding/json"
	"fmt"

	"github.com/scorredoira/dune"
)

func init() {
	dune.RegisterLib(Console, `

declare namespace console {
	export function log(...v: any[]): void
}
`)
}

var Console = []dune.NativeFunction{
	{
		Name:      "console.log",
		Arguments: -1,
		Function: func(this dune.Value, args []dune.Value, vm *dune.VM) (dune.Value, error) {
			for _, v := range args {
				var s string
				switch v.Type {
				case dune.String, dune.Int, dune.Float, dune.Bool:
					s = v.ToString()
				default:
					b, err := json.MarshalIndent(v.Export(0), "", "    ")
					if err != nil {
						return dune.NullValue, err
					}
					s = string(b)
				}

				fmt.Fprintln(vm.GetStdout(), s)
			}
			return dune.NullValue, nil
		},
	},
}
