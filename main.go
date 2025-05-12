package main

import (
	"fmt"
	"os"
	"paas/pkg/argocd"
	"paas/pkg/datadog"
	"time"

	"github.com/spf13/cobra"
)

const (
	downtimeDuration = 5 * time.Minute
)

var (
	appName string

	syncCommand = &cobra.Command{
		Use:   "sync",
		Short: "sync an application with argocd",
		RunE: func(cmd *cobra.Command, args []string) error {
			app := datadog.NewDatadogApp(datadog.WithName(appName), datadog.WithStartDate(time.Now()), datadog.WithEndDate(time.Now().Add(downtimeDuration)))

			if err := app.ActivateDowntime(); err != nil {
				return err
			}

			argoApp := argocd.NewArgocdAppHandler()

			if err := argoApp.Sync(&argocd.ArgocdAppReq{Name: appName}); err != nil {
				return fmt.Errorf("error while syncing argocd %s app, %w", appName, err)
			}

			if err := app.RemoveDownTime(); err != nil {
				return err
			}

			return nil
		},
	}
)

func initSyncCommand() *cobra.Command {
	syncCommand.Flags().StringVarP(&appName, "app", "a", "", "name of the application to sync")
	syncCommand.MarkFlagRequired("app")

	return syncCommand
}

func main() {
	cmd := initSyncCommand()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
