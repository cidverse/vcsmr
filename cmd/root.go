package cmd

import (
	"os"
	"strings"

	"github.com/cidverse/cidverseutils/zerologconfig"
	"github.com/spf13/cobra"
)

var cfg zerologconfig.LogConfig
var configFile string

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  `vcsmr`,
		Long: ``,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			zerologconfig.Configure(cfg)
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
			os.Exit(0)
		},
	}

	cmd.PersistentFlags().StringVar(&cfg.LogLevel, "log-level", "info", "log level - allowed: "+strings.Join(zerologconfig.ValidLogLevels, ","))
	cmd.PersistentFlags().StringVar(&cfg.LogFormat, "log-format", "color", "log format - allowed: "+strings.Join(zerologconfig.ValidLogFormats, ","))
	cmd.PersistentFlags().BoolVar(&cfg.LogCaller, "log-caller", false, "include caller in log functions")
	cmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")

	cmd.AddCommand(listCmd())
	cmd.AddCommand(reviewCmd())
	cmd.AddCommand(versionCmd())

	return cmd
}

// Execute executes the root command.
func Execute() error {
	return rootCmd().Execute()
}
