package cmd

import (
	"fmt"
	"github.com/jlewi/ddctl/pkg/ddog"
	"io"
	"os"
	"path/filepath"

	"github.com/go-logr/zapr"
	"github.com/jlewi/monogo/yamlfiles"
	"go.uber.org/zap"
	yaml "sigs.k8s.io/yaml/goyaml.v3"

	"github.com/jlewi/monogo/helpers"

	"github.com/jlewi/ddctl/api"
	"github.com/jlewi/ddctl/pkg/application"
	"github.com/jlewi/ddctl/pkg/version"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewLinksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "links",
	}
	cmd.AddCommand(NewBuildURL())
	cmd.AddCommand(NewParseURL())
	return cmd
}

// NewBuildURL creates a command to turn queries into URLs
func NewBuildURL() *cobra.Command {
	var patchFile string
	var open bool
	cmd := &cobra.Command{
		Use: "build",
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {
				app := application.NewApp()
				if err := app.LoadConfig(cmd); err != nil {
					return err
				}
				if err := app.SetupLogging(); err != nil {
					return err
				}

				version.LogVersion()

				nodes, err := yamlfiles.Read(patchFile)
				if err != nil {
					return errors.Wrapf(err, "Error reading file %v", patchFile)
				}

				log := zapr.NewLogger(zap.L())
				for _, n := range nodes {
					link := &api.DatadogLink{}
					switch n.GetKind() {
					case api.LinkGVK.Kind:
						if err := n.YNode().Decode(link); err != nil {
							return errors.Wrapf(err, "Error decoding %v", n.GetKind())
						}
					default:
						log.Info("Skipping object of unknown kind", "kind", n.GetKind(), "name", n.GetName(), "knownKinds", []string{api.LinkGVK.Kind})
					}
					u, err := ddog.BuildURL(link)
					if err != nil {
						return err
					}
					fmt.Printf("Datadog URL:\n%v\n", u)
					if open {
						if err := browser.OpenURL(u); err != nil {
							return errors.Wrapf(err, "Error opening URL %v", u)
						}
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

	cmd.Flags().StringVarP(&patchFile, "--filename", "f", "", "A file containing the YAML object containing the link.")
	cmd.Flags().BoolVarP(&open, "open", "", false, "Open the URL in a browser")
	return cmd
}

// NewParseURL creates a command to parse URLs
func NewParseURL() *cobra.Command {
	var panesFile string
	var logUrl string
	var name string
	cmd := &cobra.Command{
		Use: "parse",
		Run: func(cmd *cobra.Command, args []string) {
			err := func() error {
				app := application.NewApp()
				if err := app.LoadConfig(cmd); err != nil {
					return err
				}
				if err := app.SetupLogging(); err != nil {
					return err
				}

				version.LogVersion()

				var o io.Writer

				if panesFile != "" {
					f, err := os.Create(panesFile)
					if err != nil {
						return errors.Wrapf(err, "Error creating file %v", panesFile)
					}
					defer f.Close()

					if name == "" {
						// Default to the name of the file
						filename := filepath.Base(panesFile)

						// Strip the suffix (file extension)
						name = filename[:len(filename)-len(filepath.Ext(filename))]
					}

					o = f
				} else {
					o = os.Stdout
				}

				link, err := ddog.URLToLink(logUrl)
				if err != nil {
					return errors.Wrapf(err, "Error parsing URL")
				}
				link.Metadata.Name = name

				// Pretty print the json of the panes to the file
				encoder := yaml.NewEncoder(o)
				encoder.SetIndent(2)
				if err := encoder.Encode(link); err != nil {
					return errors.Wrapf(err, "Error writing Link to file")
				}
				return nil
			}()

			if err != nil {
				fmt.Printf("Error running request;\n %+v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&panesFile, "link-file", "o", "", "File to write the yaml to. If not specified the Link will be written to stdout.")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Name to give the resource when saving to a file")
	cmd.Flags().StringVarP(&logUrl, "url", "u", "", "The URL to parse")
	helpers.IgnoreError(cmd.MarkFlagRequired("url"))
	return cmd
}
