package moonshot

import (
	"context"
	"fmt"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"
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
	root.PersistentFlags().StringP("config", "c", "", "Config file path override")

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
