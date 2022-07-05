# 🌔 Moonshot 

A boilerplate Go library for quickly setting up your next moonshot idea!

## Features

* Config management
    * Create struct, pass its pointer to `moonshot.App`. 
    * Moonshot will take care of loading configs from environment/files. 
    * File can be overriden by `--config` flag also.
    * You can run `./myapp configs` to see the actual loaded configs.
* HTTP Server setup
    * HTTP server is pre-configured with graceful shutdown enabled.
    * Server is pre-configured with handlers for `/health`, NotFound, MethodNotAllowed.
    * Panic recovery is enabled.
    * You can set the `Routes` field in `moonshot.App` to add custom routes or override.
* Errors package
    * An easy-to-use errors package with common category of errors pre-defined.
    * Just do `errors.ErrInvalid.WithMsgf()` or `WithCausef()` to add additional context.
* Logging
    * `log` package is automatically configured based on `--log-level` and `--log-format` flags.
    * Pass log-context using `log.Inject(ctx, fields)` 

## Usage

1. Create `main.go`.
2. Initiailise `moonshot.App`:

   ```go
   package main
   
   import "github.com/spy16/moonshot"

   var myConfig struct {
       Database string `mapstructure:"database"`
   }

   func main() {
    app := moonshot.App{
        Name: "myapp",
        Short: "MyApp does cool things",
        CfgPtr: &myConfig,
        Routes: func(r *chi.Mux) {
            r.Get("/", myAppHomePageHandler)
        },
    }

    os.Exit(app.Launch())
   }
   ```
3. Build the app `go build -o myapp main.go`
4. Run the app: 

   ```shell
   $ ./myapp --help
   MyApp does cool things
   
   Usage:
   myapp [command]
   
   Available Commands:
   completion  Generate the autocompletion script for the specified shell
   configs     Show currently loaded configurations
   help        Help about any command
   serve       Start HTTP server.
   
   Flags:
   -c, --config string   Config file path override
   -h, --help            help for moonshot-demo
   
   Use "moonshot-demo [command] --help" for more information about a command.
   ```

* You can run `./myapp serve --addr="localhost:8080"` for starting server.
* You can pass `--static-dir` and `--static-route` flags to `serve` command for serving static files.

> **Note**: Refer `./_example` for a demo application.
