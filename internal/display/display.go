package display

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/karrick/gows"
	"github.com/spf13/viper"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/queryresult"
)

// ShowQueryOutput displays the output using the proper formatter as applicable
func ShowQueryOutput(ctx context.Context, result *queryresult.Result) int {
	rowErrors := 0

	switch cmdconfig.Viper().GetString(constants.ArgOutput) {
	case constants.OutputFormatJSON:
		rowErrors = displayJSON(ctx, result)
	case constants.OutputFormatCSV:
		rowErrors = displayCSV(ctx, result)
	case constants.OutputFormatLine:
		rowErrors = displayLine(ctx, result)
	case constants.OutputFormatTable:
		rowErrors = displayTable(ctx, result)
	}

	if shouldShowQueryTiming() {
		PrintTiming(result.Timing)
	}
	// return the number of rows that returned errors
	return rowErrors
}

type ShowWrappedTableOptions struct {
	AutoMerge        bool
	HideEmptyColumns bool
	Truncate         bool
}

func ShowWrappedTable(headers []string, rows [][]string, opts *ShowWrappedTableOptions) {
	if opts == nil {
		opts = &ShowWrappedTableOptions{}
	}
	t := table.NewWriter()

	t.SetStyle(table.StyleDefault)
	t.Style().Format.Header = text.FormatDefault
	t.SetOutputMirror(os.Stdout)

	rowConfig := table.RowConfig{AutoMerge: opts.AutoMerge}
	colConfigs, headerRow := getColumnSettings(headers, rows, opts)

	t.SetColumnConfigs(colConfigs)
	t.AppendHeader(headerRow)

	for _, row := range rows {
		rowObj := table.Row{}
		for _, col := range row {
			rowObj = append(rowObj, col)
		}
		t.AppendRow(rowObj, rowConfig)
	}
	t.Render()
}

func GetMaxCols() int {
	colsAvailable, _, _ := gows.GetWinSize()
	// check if STEAMPIPE_DISPLAY_WIDTH env variable is set
	if viper.IsSet(constants.ArgDisplayWidth) {
		colsAvailable = viper.GetInt(constants.ArgDisplayWidth)
	}
	return colsAvailable
}

// calculate and returns column configuration based on header and row content
func getColumnSettings(headers []string, rows [][]string, opts *ShowWrappedTableOptions) ([]table.ColumnConfig, table.Row) {
	colConfigs := make([]table.ColumnConfig, len(headers))
	headerRow := make(table.Row, len(headers))

	sumOfAllCols := 0

	// account for the spaces around the value of a column and separators
	spaceAccounting := ((len(headers) * 3) + 1)

	for idx, colName := range headers {
		headerRow[idx] = colName

		// get the maximum len of strings in this column
		maxLen := getTerminalColumnsRequiredForString(colName)
		colHasValue := false
		for _, row := range rows {
			colVal := row[idx]
			if !colHasValue && len(colVal) > 0 {
				// the !colHasValue is necessary in the condition,
				// otherwise, even after being set, we will keep
				// evaluating the length
				colHasValue = true
			}

			// get the maximum line length of the value
			colLen := getTerminalColumnsRequiredForString(colVal)
			if colLen > maxLen {
				maxLen = colLen
			}
		}
		colConfigs[idx] = table.ColumnConfig{
			Name:     colName,
			Number:   idx + 1,
			WidthMax: maxLen,
			WidthMin: maxLen,
		}
		if opts.HideEmptyColumns && !colHasValue {
			colConfigs[idx].Hidden = true
		}
		sumOfAllCols += maxLen
	}

	// now that all columns are set to the widths that they need,
	// set the last one to occupy as much as is available - no more - no less
	sumOfRest := sumOfAllCols - colConfigs[len(colConfigs)-1].WidthMax
	// get the max cols width
	maxCols := GetMaxCols()
	if sumOfAllCols > maxCols {
		colConfigs[len(colConfigs)-1].WidthMax = (maxCols - sumOfRest - spaceAccounting)
		colConfigs[len(colConfigs)-1].WidthMin = (maxCols - sumOfRest - spaceAccounting)
		if opts.Truncate {
			colConfigs[len(colConfigs)-1].WidthMaxEnforcer = helpers.TruncateString
		}
	}

	return colConfigs, headerRow
}

// getTerminalColumnsRequiredForString returns the length of the longest line in the string
func getTerminalColumnsRequiredForString(str string) int {
	colsRequired := 0
	scanner := bufio.NewScanner(bytes.NewBufferString(str))
	for scanner.Scan() {
		line := scanner.Text()
		runeCount := utf8.RuneCountInString(line)
		if runeCount > colsRequired {
			colsRequired = runeCount
		}
	}
	return colsRequired
}

func displayLine(ctx context.Context, result *queryresult.Result) int {

	maxColNameLength, rowErrors := 0, 0
	for _, col := range result.Cols {
		thisLength := utf8.RuneCountInString(col.Name)
		if thisLength > maxColNameLength {
			maxColNameLength = thisLength
		}
	}
	itemIdx := 0

	// define a function to display each row
	rowFunc := func(row []interface{}, result *queryresult.Result) {
		recordAsString, _ := ColumnValuesAsString(row, result.Cols)
		requiredTerminalColumnsForValuesOfRecord := 0
		for _, colValue := range recordAsString {
			colRequired := getTerminalColumnsRequiredForString(colValue)
			if requiredTerminalColumnsForValuesOfRecord < colRequired {
				requiredTerminalColumnsForValuesOfRecord = colRequired
			}
		}

		lineFormat := fmt.Sprintf("%%-%ds | %%s\n", maxColNameLength)
		multiLineFormat := fmt.Sprintf("%%-%ds | %%-%ds", maxColNameLength, requiredTerminalColumnsForValuesOfRecord)

		fmt.Printf("-[ RECORD %-2d ]%s\n", (itemIdx + 1), strings.Repeat("-", 75)) //nolint:forbidigo // intentional use of fmt
		for idx, column := range recordAsString {
			lines := strings.Split(column, "\n")
			if len(lines) == 1 {
				fmt.Printf(lineFormat, result.Cols[idx].Name, lines[0]) //nolint:forbidigo // intentional use of fmt
			} else {
				for lineIdx, line := range lines {
					if lineIdx == 0 {
						// the first line
						fmt.Printf(multiLineFormat, result.Cols[idx].Name, line) //nolint:forbidigo // intentional use of fmt
					} else {
						// next lines
						fmt.Printf(multiLineFormat, "", line) //nolint:forbidigo // intentional use of fmt
					}

					// is this not the last line of value?
					if lineIdx < len(lines)-1 {
						fmt.Printf(" +\n") //nolint:forbidigo // intentional use of fmt
					} else {
						fmt.Printf("\n") //nolint:forbidigo // intentional use of fmt
					}

				}
			}
		}
		itemIdx++

	}

	// call this function for each row
	if err := iterateResults(result, rowFunc); err != nil {
		error_helpers.ShowError(ctx, err)
		rowErrors++
		return rowErrors
	}
	return rowErrors
}

type resultMetadata struct {
	RowsReturned int    `json:"rows_returned"`
	Duration     string `json:"duration_ms"`
}
type jsonOutput struct {
	Columns  []columnDef              `json:"columns"`
	Rows     []map[string]interface{} `json:"rows"`
	Metadata resultMetadata           `json:"metadata"`
}

type columnDef struct {
	Name         string `json:"name"`
	DataType     string `json:"data_type"`
	OriginalName string `json:"original_name,omitempty"`
}

type uniqueNameGenerator struct {
	lookup map[string]struct{}
}

// ctor
func newUniqueNameGenerator() *uniqueNameGenerator {
	return &uniqueNameGenerator{
		lookup: make(map[string]struct{}),
	}
}

// getUniqueName returns a unique name based on the input name
func (g *uniqueNameGenerator) getUniqueName(name string) string {
	// ensure a unique column name
	for {
		// check the lookup to see if this name exists
		if _, exists := g.lookup[name]; !exists {
			// name is unique - we are done
			break
		}
		// name is not unique - generate a new name
		// store the original name
		originalName := name

		// generate a new name
		name = fmt.Sprintf("%s_%s", originalName, utils.RandomString(4))
	}
	// add the unique name into the lookup
	g.lookup[name] = struct{}{}
	return name
}

func displayJSON(ctx context.Context, result *queryresult.Result) int {
	rowErrors := 0
	var op = jsonOutput{
		Metadata: resultMetadata{
			Duration: getDurationString(result.Timing.Duration),
		},
	}

	// create a unique name generator
	nameGenerator := newUniqueNameGenerator()

	// add column defs to the JSON output
	for _, col := range result.Cols {
		c := columnDef{
			Name:     nameGenerator.getUniqueName(col.Name),
			DataType: strings.ToLower(col.DataType),
		}

		// if the column name has changed, store the original
		if c.Name != col.Name {
			c.OriginalName = col.Name
		}
		// add to the column def array
		op.Columns = append(op.Columns, c)
	}

	// Define function to add each row to the JSON output
	rowFunc := func(row []interface{}, result *queryresult.Result) {
		record := map[string]interface{}{}
		for idx, col := range result.Cols {
			value, _ := ParseJSONOutputColumnValue(row[idx], col)
			// get the column def
			c := op.Columns[idx]
			// add the value under the unique column name
			record[c.Name] = value
		}

		op.Rows = append(op.Rows, record)
	}

	// call this function for each row
	if err := iterateResults(result, rowFunc); err != nil {
		error_helpers.ShowError(ctx, err)
		rowErrors++
		return rowErrors
	}
	op.Metadata.RowsReturned = len(op.Rows)

	// display the JSON
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", " ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(op); err != nil {
		error_helpers.ShowErrorWithMessage(ctx, err, "Error displaying result as JSON")
		return 0
	}
	return rowErrors
}

func displayCSV(ctx context.Context, result *queryresult.Result) int {
	rowErrors := 0
	csvWriter := csv.NewWriter(os.Stdout)
	csvWriter.Comma = []rune(cmdconfig.Viper().GetString(constants.ArgSeparator))[0]

	if cmdconfig.Viper().GetBool(constants.ArgHeader) {
		_ = csvWriter.Write(ColumnNames(result.Cols))
	}

	// print the data as it comes
	// define function display each csv row
	rowFunc := func(row []interface{}, result *queryresult.Result) {
		rowAsString, _ := ColumnValuesAsString(row, result.Cols, WithNullString(""))
		_ = csvWriter.Write(rowAsString)
	}

	// call this function for each row
	if err := iterateResults(result, rowFunc); err != nil {
		error_helpers.ShowError(ctx, err)
		rowErrors++
		return rowErrors
	}

	csvWriter.Flush()
	if csvWriter.Error() != nil {
		error_helpers.ShowErrorWithMessage(ctx, csvWriter.Error(), "unable to print csv")
	}
	return rowErrors
}

func displayTable(ctx context.Context, result *queryresult.Result) int {
	rowErrors := 0
	// the buffer to put the output data in
	outbuf := bytes.NewBufferString("")

	// the table
	t := table.NewWriter()
	t.SetOutputMirror(outbuf)
	t.SetStyle(table.StyleDefault)
	t.Style().Format.Header = text.FormatDefault

	colConfigs := []table.ColumnConfig{}
	headers := make(table.Row, len(result.Cols))

	for idx, column := range result.Cols {
		headers[idx] = column.Name
		colConfigs = append(colConfigs, table.ColumnConfig{
			Name:     column.Name,
			Number:   idx + 1,
			WidthMax: constants.MaxColumnWidth,
		})
	}

	t.SetColumnConfigs(colConfigs)
	if viper.GetBool(constants.ArgHeader) {
		t.AppendHeader(headers)
	}

	// define a function to execute for each row
	rowFunc := func(row []interface{}, result *queryresult.Result) {
		rowAsString, _ := ColumnValuesAsString(row, result.Cols)
		rowObj := table.Row{}
		for _, col := range rowAsString {
			// trim out non-displayable code-points in string
			// exfept white-spaces
			col = strings.Map(func(r rune) rune {
				if unicode.IsSpace(r) || unicode.IsGraphic(r) {
					// return if this is a white space character
					return r
				}
				return -1
			}, col)
			rowObj = append(rowObj, col)
		}
		t.AppendRow(rowObj)
	}

	// iterate each row, adding each to the table
	err := iterateResults(result, rowFunc)
	if err != nil {
		// display the error
		fmt.Println() //nolint:forbidigo // intentional use of fmt
		error_helpers.ShowError(ctx, err)
		rowErrors++
		fmt.Println() //nolint:forbidigo // intentional use of fmt
	}
	// write out the table to the buffer
	t.Render()

	// page out the table
	ShowPaged(ctx, outbuf.String())
	return rowErrors
}

type displayResultsFunc func(row []interface{}, result *queryresult.Result)

// call func displayResult for each row of results
func iterateResults(result *queryresult.Result, displayResult displayResultsFunc) error {
	for row := range *result.RowChan {
		if row == nil {
			return nil
		}
		if row.Error != nil {
			return row.Error
		}
		displayResult(row.Data, result)
	}
	// we will not get here
	return nil
}

func shouldShowQueryTiming() bool {
	outputFormat := viper.GetString(constants.ArgOutput)
	return viper.GetBool(constants.ArgTiming) && outputFormat == constants.OutputFormatTable
}

func PrintTiming(timingMetadata *queryresult.TimingMetadata) {
	durationString := getDurationString(timingMetadata.Duration)
	fmt.Printf("\nTime: %s\n", durationString) //nolint:forbidigo // intentional use of fmt
}

func getDurationString(duration time.Duration) string {
	// Calculate duration since startTime and round down to the nearest millisecond
	durationInMS := duration / time.Millisecond
	//nolint:durationcheck // we want to print the duration in milliseconds
	duration = durationInMS * time.Millisecond

	durationString := duration.String()
	return durationString
}
