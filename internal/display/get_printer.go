package display

import (
	"github.com/spf13/cobra"
	"github.com/turbot/pipe-fittings/printers"
)

// TODO is this even needed
func GetPrinter[T any](cmd *cobra.Command) (printers.ResourcePrinter[T], error) {
	return printers.GetPrinter[T](cmd)
}
