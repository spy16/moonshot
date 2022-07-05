package moonshot

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/spy16/moonshot/config"
	"github.com/spy16/moonshot/log"
)

func (cli *App) cmdShowConfigs(ctx context.Context) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "configs",
		Short: "Show currently loaded configurations",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cli.loadConfigs(cmd); err != nil {
				log.Fatalf(ctx, "failed to load configurations: %v", err)
			}

			var err error
			if format == "json" {
				err = json.NewEncoder(os.Stdout).Encode(cli.CfgPtr)
			} else {
				enc := yaml.NewEncoder(os.Stdout)
				err = enc.Encode(cli.CfgPtr)
			}

			if err != nil {
				log.Fatalf(ctx, "failed to display configs: %v", err)
			}
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "yaml", "Output format")

	return cmd
}

func (cli *App) loadConfigs(cmd *cobra.Command) error {
	opts := []config.Option{
		config.WithName(strings.ToLower(cli.Name)),
		config.WithEnv(""),
	}

	cfgFile, err := cmd.Flags().GetString("config")
	if err == nil && cfgFile != "" {
		opts = append(opts, config.WithFile(cfgFile))
	}

	if err := config.Load(cli.CfgPtr, opts...); err != nil {
		return err
	}

	return nil
}
