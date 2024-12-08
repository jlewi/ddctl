package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/go-logr/zapr"
	"github.com/jlewi/ddctl/pkg/application"
	"github.com/jlewi/ddctl/pkg/config"
	"github.com/jlewi/ddctl/pkg/ddog"
	"github.com/jlewi/ddctl/pkg/version"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// NewQueryToURL creates a command to turn queries into URLs
func NewQueryToURL() *cobra.Command {
	var query string
	var queryFile string
	var baseURL string
	var open bool
	var duration time.Duration
	var endTime string
	var layout string
	cmd := &cobra.Command{
		Use: "querytourl",
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {
				app := application.NewApp()
				if err := app.LoadConfig(cmd); err != nil {
					return err
				}
				if err := app.SetupLogging(); err != nil {
					return err
				}

				log := zapr.NewLogger(zap.L())
				version.LogVersion()

				if (query == "" && queryFile == "") || (query != "" && queryFile != "") {
					return errors.New("Exactly one of --query and --query-file must be specified")
				}

				if queryFile != "" {
					data, err := os.ReadFile(queryFile)
					if err != nil {
						return errors.Wrapf(err, "Error reading query file %v", queryFile)
					}
					query = string(data)
				}

				// Treat the query as the YAML representation of a map
				queryArgs := map[string]string{}

				if err := yaml.Unmarshal([]byte(query), &queryArgs); err != nil {
					log.Error(err, "Error unmarshalling query", "query", query)
					return errors.Wrapf(err, "Error unmarshalling query")
				}

				if app.Config.GetBaseURL() == "" {
					return errors.New("baseURL must be specified either in config.yaml or via the --base-url flag")
				}

				timeRange, err := ddog.TimeAndDurationToRange(endTime, layout, duration)
				if err != nil {
					return err
				}

				for k, v := range timeRange {
					queryArgs[k] = v
				}

				// Live controls whether logs are centered at the current time. We generally want to set that
				// to false because we will specify a particular time range in order to make the links permalinks
				// by default, however user can override this by explicitly setting live.
				if _, ok := queryArgs["live"]; !ok {
					queryArgs["live"] = "false"
				}
				u := ddog.GetLogsLink(app.Config.GetBaseURL(), queryArgs)

				fmt.Printf("Datadog URL:\n%v\n", u)
				if open {
					if err := browser.OpenURL(u); err != nil {
						return errors.Wrapf(err, "Error opening URL %v", u)
					}
				}
				return nil
			}()

			if err != nil {
				fmt.Printf("Error running request;\n %+v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&query, "query", "", "", "The datadog query")
	cmd.Flags().StringVarP(&queryFile, "query-file", "", "", "A file containing the honeycomb query")
	cmd.Flags().StringVarP(&baseURL, config.BaseURLFlagName, "", "", "The base URL for your Datadog URLs. It should be something like https://acme.datadoghq.com")
	cmd.Flags().BoolVarP(&open, "open", "", false, "Open the URL in a browser")
	cmd.Flags().DurationVarP(&duration, "duration", "d", 24*time.Hour, "The duration for the query")
	cmd.Flags().StringVarP(&endTime, "end-time", "t", "", "The end time for the query. Defaults to now if not specified.")
	cmd.Flags().StringVarP(&layout, "layout", "l", "2006-01-02 15:04 MST", "Layout for parsing time strings")
	return cmd
}
