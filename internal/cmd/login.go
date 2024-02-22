package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/cloud"
	"github.com/turbot/pipe-fittings/cmdconfig"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/error_helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

func loginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "login",
		TraverseChildren: true,
		Args:             cobra.NoArgs,
		Run:              runLoginCmd,
		Short:            "Login to Turbot Pipes",
		Long: `The Powerpipe CLI can interact with Turbot Pipes to run pipelines in a remote cloud instance. 
		
These capabilities require authenticating to Turbot Pipes. The powerpipe login command launches an interactive process for logging in 
and obtaining a temporary (30 day) token. The token is written to ~/.powerpipe/internal/{cloud host}.pptt.`,
	}

	cmdconfig.OnCmd(cmd).
		AddCloudFlags().
		AddBoolFlag(constants.ArgHelp, false, "Help for dashboard", cmdconfig.FlagOptions.WithShortHand("h"))

	return cmd
}

func runLoginCmd(cmd *cobra.Command, _ []string) {
	ctx := cmd.Context()

	log.Printf("[TRACE] login, cloud host %s", viper.Get(constants.ArgPipesHost))
	log.Printf("[TRACE] opening login web page")
	// start login flow - this will open a web page prompting user to login, and will give the user a code to enter
	var id, err = cloud.WebLogin(ctx)
	if err != nil {
		error_helpers.ShowError(ctx, err)
		exitCode = constants.ExitCodeLoginCloudConnectionFailed
		return
	}

	token, err := getToken(ctx, id)
	if err != nil {
		error_helpers.ShowError(ctx, err)
		exitCode = constants.ExitCodeLoginCloudConnectionFailed
		return
	}

	// save token
	err = cloud.SaveToken(token)
	if err != nil {
		error_helpers.ShowError(ctx, err)
		exitCode = constants.ExitCodeLoginCloudConnectionFailed
		return
	}

	displayLoginMessage(ctx, token)
}

func getToken(ctx context.Context, id string) (loginToken string, err error) {
	log.Printf("[TRACE] prompt for verification code")

	//nolint:forbidigo // intentional
	fmt.Println()
	retries := 0
	for {
		var code string
		code, err = promptUserForString("Enter verification code: ")
		error_helpers.FailOnError(err)
		if code != "" {
			log.Printf("[TRACE] get login token")
			// use this code to get a login token and store it
			loginToken, err = cloud.GetLoginToken(ctx, id, code)
			if err == nil {
				return loginToken, nil
			}
		}
		if err != nil {
			// a code was entered but it failed - inc retry count
			log.Printf("[TRACE] GetLoginToken failed with %s", err.Error())
		}
		retries++

		// if we have used our retries, break out before displaying wanring - we will display an error
		if retries == 3 {
			return "", sperr.New("Too many attempts.")
		}

		if err != nil {
			error_helpers.ShowWarning(err.Error())
		}
		log.Printf("[TRACE] Retrying")
	}
}

func displayLoginMessage(ctx context.Context, token string) {
	userName, err := cloud.GetUserName(ctx, token)
	error_helpers.FailOnError(sperr.WrapWithMessage(err, "failed to read user name"))

	//nolint:forbidigo // intentional
	fmt.Printf("\nLogged in as: %s\n\n", constants.Bold(userName))
}

func promptUserForString(prompt string) (string, error) {
	//nolint:forbidigo // intentional
	fmt.Print(prompt)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		// handle ctrl+d
		//nolint:forbidigo // intentional
		fmt.Println()
		os.Exit(0)
	}

	err := scanner.Err()
	if err != nil {
		return "", sperr.Wrap(err)
	}
	code := scanner.Text()

	return code, nil
}
