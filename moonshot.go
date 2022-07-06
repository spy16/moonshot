package moonshot

import (
	"context"
	"fmt"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
	"github.com/spy16/moonshot/log"
)

// App represents an instance of app command. Invoke App.Launch()
// in `main()`.
type App struct {
	Name   string
	Short  string
	Long   string
	CfgPtr interface{}
	Routes func(r *chi.Mux)
}

func (app *App) Launch(ctx context.Context, cmds ...*cobra.Command) int {
	root := &cobra.Command{
		Use:   fmt.Sprintf("%s <command> [flags] <args>", app.Name),
		Short: app.Short,
		Long:  app.Long,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	flags := root.PersistentFlags()

	var logLevel, logFormat string
	flags.StringP("config", "c", "", "Config file path override")
	flags.StringVar(&logLevel, "log-level", "info", "Log level")
	flags.StringVar(&logFormat, "log-format", "text", "Log format (json/text)")

	root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		log.Setup(logLevel, logFormat)
	}

	root.AddCommand(cmds...)
	root.AddCommand(
		app.cmdServe(ctx),
		app.cmdShowConfigs(ctx),
	)

	if err := root.Execute(); err != nil {
		return 1
	}
	return 0
}
