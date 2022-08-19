package moonshot

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/spf13/cobra"

	"github.com/spy16/moonshot/errors"
	"github.com/spy16/moonshot/httputils"
	"github.com/spy16/moonshot/log"
)

func (app *App) cmdServe(ctx context.Context) *cobra.Command {
	var graceDur time.Duration
	var addr, staticDir, staticRoute string
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HTTP server.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := app.loadConfigs(cmd); err != nil {
				log.Fatalf(ctx, "failed to load configs: %v", err)
				return
			}

			router := chi.NewRouter()
			router.NotFound(notFoundHandler())
			router.MethodNotAllowed(methodNotAllowedHandler())
			router.Get("/health", pingHandler(map[string]interface{}{
				"status": "ok",
			}))

			if err := app.Routes(router); err != nil {
				log.Fatalf(ctx, "route setup failed: %v", err)
			}

			if app.StaticFS != nil {
				router.Mount(staticRoute, staticHandler(app.StaticFS))
			} else if staticDir != "" {
				router.Mount(staticRoute, http.StripPrefix(staticRoute, http.FileServer(http.Dir(staticDir))))
			}

			log.Infof(ctx, "starting server at '%s'...", addr)
			if err := httputils.GracefulServe(ctx, graceDur, addr, router); err != nil {
				log.Fatalf(ctx, "server exited with error: %v", err)
			}
		},
	}

	cmd.Flags().StringVarP(&addr, "addr", "a", ":8080", "Bind address for HTTP server")
	cmd.Flags().StringVarP(&staticDir, "static-dir", "D", "", "Directory to serve static files from")
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

func staticHandler(staticFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		path := filepath.Clean(r.URL.Path)
		if path == "/" { // Add other paths that you route on the UI side here
			path = "index.html"
		}
		path = strings.TrimPrefix(path, "/")

		file, err := staticFS.Open(path)
		if err != nil {
			if os.IsNotExist(err) {
				log.Warnf(r.Context(), "file '%s' not found: %v", path, err)
				http.NotFound(w, r)
				return
			}
			log.Warnf(r.Context(), "failed to open '%s': %v", path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(path))
		w.Header().Set("Content-Type", contentType)
		if strings.HasPrefix(path, "static/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		}

		stat, err := file.Stat()
		if err == nil && stat.Size() > 0 {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
		}

		n, _ := io.Copy(w, file)
		log.Debugf(r.Context(), "%d bytes copied to client from file '%s'", n, path)
	}
}
