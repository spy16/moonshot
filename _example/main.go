package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"

	"github.com/spy16/moonshot"
	"github.com/spy16/moonshot/pkg/httputils"
)

var appCfg struct {
	Addr      string `mapstructure:"addr" yaml:"addr" json:"addr"`
	LogLevel  string `mapstructure:"log_level" yaml:"log_level" json:"log_level"`
	LogFormat string `mapstructure:"log_format" yaml:"log_format" json:"log_format"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app := &moonshot.App{
		Name:   "moonshot-demo",
		Short:  "A sample moonshot app setup",
		CfgPtr: &appCfg,
		Routes: func(r *chi.Mux) error {
			// set up any custom routes here.
			r.Get("/hello", func(wr http.ResponseWriter, req *http.Request) {
				httputils.Respond(wr, req, http.StatusOK, "Hello!")
			})

			return nil
		},
	}

	os.Exit(app.Launch(ctx))
}
