package initialisation

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
)

type InitResult struct {
	error_helpers.ErrorAndWarnings
	Messages []string

	// allow overriding of the display functions
	DisplayMessage func(ctx context.Context, m string)
	DisplayWarning func(ctx context.Context, w string)
}

func (r *InitResult) AddMessage(messages ...string) {
	r.Messages = append(r.Messages, messages...)
}

func (r *InitResult) AddWarnings(warnings ...string) {
	r.Warnings = append(r.Warnings, warnings...)
}

func (r *InitResult) HasMessages() bool {
	return len(r.Warnings)+len(r.Messages) > 0
}

func (r *InitResult) DisplayMessages() {
	if r.DisplayMessage == nil {
		r.DisplayMessage = func(ctx context.Context, m string) {
			fmt.Println(m) //nolint:forbidigo // TODO
		}
	}
	if r.DisplayWarning == nil {
		r.DisplayWarning = func(ctx context.Context, w string) {
			error_helpers.ShowWarning(w)
		}
	}
	for _, w := range r.Warnings {
		r.DisplayWarning(context.Background(), w)
	}
	// do not display message in json or csv output mode
	output := viper.Get(constants.ArgOutput)
	if output == constants.OutputFormatJSON || output == constants.OutputFormatCSV {
		return
	}
	for _, m := range r.Messages {
		r.DisplayMessage(context.Background(), m)
	}
}

func (r *InitResult) Merge(other InitResult) {
	r.ErrorAndWarnings.Merge(other.ErrorAndWarnings)

	r.AddMessage(other.Messages...)
}
