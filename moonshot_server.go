package moonshot

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"

	"github.com/spy16/moonshot/errors"
	"github.com/spy16/moonshot/httputils"
	"github.com/spy16/moonshot/log"
)

func (cli *App) cmdServe(ctx context.Context) *cobra.Command {
	var graceDur time.Duration
	var addr, staticDir, staticRoute string
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HTTP server.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cli.loadConfigs(cmd); err != nil {
				log.Fatalf(ctx, "failed to load configs: %v", err)
				return
			}

			router := chi.NewRouter()
			router.NotFound(notFoundHandler())
			router.MethodNotAllowed(methodNotAllowedHandler())
			router.Get("/health", pingHandler(map[string]interface{}{
				"status": "ok",
			}))
			if staticRoute != "" {
				router.Mount(staticRoute, http.StripPrefix(staticRoute, http.FileServer(http.Dir(staticDir))))
			}

			cli.Routes(router)

			log.Infof(ctx, "starting server at '%s'...", addr)
			if err := httputils.GracefulServe(ctx, graceDur, addr, router); err != nil {
				log.Fatalf(ctx, "server exited with error: %v", err)
			}
		},
	}

	cmd.Flags().StringVarP(&addr, "addr", "a", ":8080", "Bind address for HTTP server")
	cmd.Flags().StringVarP(&staticDir, "static-dir", "D", "./app", "Directory to serve static files from")
	cmd.Flags().StringVarP(&staticRoute, "static-route", "R", "/", "Route to serve static files under")
	cmd.Flags().DurationVarP(&graceDur, "grace-period", "G", 5*time.Second, "Grace period for shutdown")
	return cmd
}

func methodNotAllowedHandler() http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		httputils.Respond(wr, req, http.StatusMethodNotAllowed,
			errors.ErrInvalid.WithMsgf("%s not allowed for %s", req.Method, req.URL.Path))
	}
}

func notFoundHandler() http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		httputils.Respond(wr, req, http.StatusNotFound,
			errors.ErrNotFound.WithMsgf("endpoint '%s %s' not found", req.Method, req.URL.Path))
	}
}

func pingHandler(info map[string]interface{}) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		httputils.Respond(wr, req, http.StatusOK, info)
	}
}
