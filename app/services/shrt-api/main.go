package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/yashshah7197/shrt/app/services/shrt-api/handlers"
	"github.com/yashshah7197/shrt/business/sys/auth"
	"github.com/yashshah7197/shrt/foundation/keystore"

	"github.com/ardanlabs/conf"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var build = "develop"

func main() {
	// Construct the application logger.
	logger, err := initLogger("SHRT-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(logger); err != nil {
		logger.Errorw("startup", "ERROR", err)
		os.Exit(1)
	}
}

func run(logger *zap.SugaredLogger) error {
	// =============================================================================================
	// GOMAXPROCS
	// =============================================================================================

	// Set the correct number of threads for the service based on what is
	// available either by the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	logger.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =============================================================================================
	// Configuration
	// =============================================================================================

	cfg := struct {
		conf.Version
		Web struct {
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
		}
		Auth struct {
			KeysFolder  string `conf:"default:zarf/keys/"`
			ActiveKeyID string `conf:"default:ecdf8542-fbf3-404d-acdc-f41527a0c3c8"`
		}
	}{
		Version: conf.Version{
			SVN:  build,
			Desc: "Copyright Yash Shah, 2022",
		},
	}

	const prefix = "SHRT"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		return fmt.Errorf("parsing config: %w", err)
	}

	// =============================================================================================
	// Application Starting
	// =============================================================================================

	logger.Infow("starting service", "version", build)
	defer logger.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	logger.Infow("startup", "config", out)

	expvar.NewString("build").Set(build)

	// =============================================================================================
	// Initialize Authentication & Authorization Support
	// =============================================================================================
	logger.Infow("startup", "status", "initializing authentication & authorization support")

	// Construct a keystore based on the key files stored in the specified directory
	ks, err := keystore.NewFS(os.DirFS(cfg.Auth.KeysFolder))
	if err != nil {
		return fmt.Errorf("reading keys from keys folder: %w", err)
	}

	auth, err := auth.New(cfg.Auth.ActiveKeyID, ks)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	// =============================================================================================
	// Start Debug Service
	// =============================================================================================

	logger.Infow("startup", "status", "debug router started", "host", cfg.Web.DebugHost)

	// The Debug function returns a mux to listen and serve on for all the debug related endpoints.
	// This includes the standard library endpoints.

	// Construct the mux for the debug calls.
	debugMux := handlers.DebugMux(build, logger)

	// Start the service listening for debug requests.
	// Not concerned with shutting this down with load shedding.
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debugMux); err != nil {
			logger.Errorw("shutdown, status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// =============================================================================================
	// Start API Service
	// =============================================================================================

	logger.Infow("startup", "status", "initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS. Use a buffered
	// channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Construct the mux for API calls.
	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown: shutdown,
		Logger:   logger,
		Auth:     auth,
	})

	// Construct a server to service requests against the mux.
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(logger.Desugar()),
	}

	// Make a channel to listen to errors coming from the listener. Use a buffered channel so the
	// goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for API requests.
	go func() {
		logger.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// =============================================================================================
	// Shutdown
	// =============================================================================================

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer logger.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func initLogger(service string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}
