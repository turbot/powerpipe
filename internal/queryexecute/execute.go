package queryexecute

import (
	"context"
	"github.com/turbot/pipe-fittings/modconfig"
	"github.com/turbot/pipe-fittings/utils"
	"github.com/turbot/powerpipe/internal/initialisation"
	"github.com/turbot/powerpipe/internal/queryresult"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
	"log"
)

// TODO NOT NEEDED

func Execute(ctx context.Context, initData *initialisation.InitData) error {
	utils.LogTime("queryexecute.Execute start")
	defer utils.LogTime("queryexecute.Execute end")

	// failures return the number of queries that failed and also the number of rows that
	// returned errors

	//t := time.Now()
	// we expect a single query only
	if len(initData.Targets) != 1 {
		return sperr.New("expected a single query to execute, got %d", len(initData.Targets))
	}

	qp, ok := initData.Targets[0].(modconfig.QueryProvider)
	if !ok {
		return sperr.New("expected a query or resource which implements QueryProvider as target,  got %T", initData.Targets[0])
	}

	// resolve query with args
	// TODO HANDLE ARGS
	resolvedQuery, err := qp.GetResolvedQuery(&modconfig.QueryArgs{})
	if err != nil {
		return err
	}

	utils.LogTime("query.execute.executeQuery start")
	defer utils.LogTime("query.execute.executeQuery end")

	utils.LogTime("db.ExecuteQuery start")
	defer utils.LogTime("db.ExecuteQuery end")

	resultsStreamer := queryresult.NewResultStreamer()
	result, err := initData.Client.Execute(ctx, resolvedQuery.ExecuteSQL, resolvedQuery.Args...)
	if err != nil {
		return err
	}
	go func() {
		resultsStreamer.StreamResult(result)
		resultsStreamer.Close()
	}()

	//rowErrors := 0 // get the number of rows that returned an error
	//// print the data as it comes
	for r := range resultsStreamer.Results {
		// todo show output
		//rowErrors = display.ShowOutput(ctx, r)
		log.Println(r)
		// signal to the resultStreamer that we are done with this result
		resultsStreamer.AllResultsRead()
	}
	return err

}
