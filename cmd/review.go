package cmd

import (
	"fmt"
	"os"

	"github.com/cidverse/go-ptr"
	"github.com/cidverse/go-rules/pkg/expr"
	"github.com/cidverse/go-vcsapp/pkg/platform/api"
	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	"github.com/cidverse/vcspr/pkg/config"
	"github.com/cidverse/vcspr/pkg/mrutil"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func reviewCmd() *cobra.Command {
	var configPath string
	var method string

	cmd := &cobra.Command{
		Use:     "review",
		Aliases: []string{},
		Short:   `automatic rebase / approve / merge based on the provided rules`,
		Run: func(cmd *cobra.Command, args []string) {
			// Load rules config
			conf, err := config.LoadConfig(configPath)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to load config")
			}

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
			for _, mr := range mrs {
				// diff
				diff, err := platform.MergeRequestDiff(mr.Repository, mr)
				if err != nil {
					log.Error().Err(err).Msg("failed to get merge request diff")
					continue
				}

				// evaluate rules
				mrContext := mrutil.GenerateMRContext(mr, diff)
				var matchedActions []string
				for _, rule := range conf.Rules {
					result, err := expr.EvalBooleanExpression(rule.Expression, mrContext)
					if err != nil {
						log.Error().Err(err).Str("expression", rule.Expression).Msg("ignoring rule due to error")
						continue
					}
					log.Debug().Str("expression", rule.Expression).Bool("result", result).Msg("evaluating rule")

					if result {
						matchedActions = append(matchedActions, rule.Action)
					}
				}

				if method == "cli" && mr.Repository.PlatformType == "gitlab" {
					for _, action := range matchedActions {
						switch action {
						case "rebase":
							fmt.Printf("glab mr rebase %d --skip-ci --repo %s\n", mr.Number, mr.Repository.Namespace+"/"+mr.Repository.Name)
						case "approve":
							fmt.Printf("glab mr approve %d --repo %s\n", mr.Number, mr.Repository.Namespace+"/"+mr.Repository.Name)
						case "merge":
							fmt.Printf("glab mr merge %d --rebase --squash --auto-merge --repo %s\n", mr.Number, mr.Repository.Namespace+"/"+mr.Repository.Name)
						case "close":
							fmt.Printf("glab mr close %d --repo %s\n", mr.Number, mr.Repository.Namespace+"/"+mr.Repository.Name)
						default:
							log.Warn().Str("action", action).Msg("unknown action in config")
						}
					}
				} else if method == "cli" && mr.Repository.PlatformType == "github" {
					for _, action := range matchedActions {
						switch action {
						case "rebase":
							fmt.Printf("gh pr rebase %d --repo %s\n", mr.Number, mr.Repository.Namespace+"/"+mr.Repository.Name)
						case "approve":
							fmt.Printf("gh pr review %d --approve --repo %s\n", mr.Number, mr.Repository.Namespace+"/"+mr.Repository.Name)
						case "merge":
							fmt.Printf("gh pr merge %d --squash --repo %s\n", mr.Number, mr.Repository.Namespace+"/"+mr.Repository.Name)
						case "close":
							fmt.Printf("gh pr close %d --repo %s\n", mr.Number, mr.Repository.Namespace+"/"+mr.Repository.Name)
						default:
							log.Warn().Str("action", action).Msg("unknown action in config")
						}
					}
				} else if method == "api" {
					for _, action := range matchedActions {
						switch action {
						case "rebase":
							// TODO: implement rebase via API
						case "approve":
							err = platform.SubmitReview(mr.Repository, mr, true, nil)
						case "merge":
							err = platform.Merge(mr.Repository, mr, api.MergeStrategyOptions{Squash: ptr.True()})
						case "close":
							// TODO: implement close via API
						default:
							log.Warn().Str("action", action).Msg("unknown action in config")
						}
					}
				} else {
					log.Error().Str("method", method).Msg("unknown method, please use 'cli' or 'api'")
					os.Exit(1)
				}
			}
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "review.yml", "path to review config")
	cmd.Flags().StringVarP(&method, "method", "m", "cli", "method to use to apply the actions (cli, api)")

	return cmd
}
