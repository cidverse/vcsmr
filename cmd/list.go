package cmd

import (
	"fmt"
	"os"

	"github.com/cidverse/cidverseutils/core/clioutputwriter"
	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{},
		Short:   `list all merge requests`,
		Run: func(cmd *cobra.Command, args []string) {
			format, _ := cmd.Flags().GetString("format")
			columns, _ := cmd.Flags().GetStringSlice("columns")

			// platform
			platform, err := vcsapp.GetPlatformFromEnvironment()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to configure platform from environment")
			}

			// query
			mrs, err := vcsapp.ListMergeRequests(platform)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to list merge requests")
				os.Exit(1)
			}

			// data
			data := clioutputwriter.TabularData{
				Headers: []string{"ID", "REPO_URL", "TITLE", "AUTHOR_ID", "AUTHOR_NAME"},
				Rows:    [][]interface{}{},
			}
			for _, mr := range mrs {
				data.Rows = append(data.Rows, []interface{}{
					mr.Id,
					mr.Repository.URL,
					mr.Title,
					mr.Author.ID,
					mr.Author.Name,
				})
			}

			// filter columns
			if len(columns) > 0 {
				data = clioutputwriter.FilterColumns(data, columns)
			}

			// print
			err = clioutputwriter.PrintData(os.Stdout, data, clioutputwriter.Format(format))
			if err != nil {
				log.Fatal().Err(err).Msg("failed to print data")
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringP("format", "f", string(clioutputwriter.DefaultOutputFormat()), fmt.Sprintf("output format %s", clioutputwriter.SupportedOutputFormats()))
	cmd.Flags().StringSliceP("columns", "c", []string{}, "columns to display")

	return cmd
}
