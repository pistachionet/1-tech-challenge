package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/navid/blog/internal/config"
	"github.com/navid/blog/internal/middleware"
	"github.com/navid/blog/internal/routes"
	"github.com/navid/blog/internal/services"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "server encountered an error: %s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// Load and validate environment config
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("[in main.run] failed to load config: %w", err)
	}

	// Create a structured logger, which will print logs in json format to the
	// writer we specify.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))

	// Create a new DB connection using environment config
	logger.DebugContext(ctx, "Connecting to database")
	db, err := sql.Open("pgx", fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUserName,
		cfg.DBUserPassword,
		cfg.DBName,
		cfg.DBPort,
	))
	if err != nil {
		return fmt.Errorf("[in main.run] failed to open database: %w", err)
	}

	// Ping the database to verify connection
	logger.DebugContext(ctx, "Pinging database")
	if err = db.PingContext(ctx); err != nil {
		return fmt.Errorf("[in main.run] failed to ping database: %w", err)
	}

	defer func() {
		logger.DebugContext(ctx, "Closing database connection")
		if err = db.Close(); err != nil {
			logger.ErrorContext(ctx, "Failed to close database connection", "err", err)
		}
	}()

	logger.InfoContext(ctx, "Connected successfully to the database")

	// Create a new users service
	usersService := services.NewUsersService(logger, db)

	// Create a serve mux to act as our route multiplexer
	mux := http.NewServeMux()

	// Add our routes to the mux
	routes.AddRoutes(
		mux,
		logger,
		usersService,
		fmt.Sprintf("http://%s:%s", cfg.Host, cfg.Port),
	)

	// Wrap the mux with middleware
	wrappedMux := middleware.Logger(logger)(mux)
	wrappedMux = middleware.Recovery(logger)(wrappedMux)

	// Create a new http server with our mux as the handler
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: wrappedMux,
	}

	errChan := make(chan error)

	// Server run context
	ctx, done := context.WithCancel(ctx)
	defer done()

	// Handle graceful shutdown with go routine on SIGINT
	go func() {
		// create a channel to listen for SIGINT and then block until it is received
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig

		logger.DebugContext(ctx, "Received SIGINT, shutting down server")

		// Create a context with a timeout to allow the server to shut down gracefully
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		// Shutdown the server. If an error occurs, send it to the error channel
		if err = httpServer.Shutdown(ctx); err != nil {
			errChan <- fmt.Errorf("[in main.run] failed to shutdown http server: %w", err)
			return
		}

		// Close the idle connections channel, unblocking `run()`
		done()
	}()

	// Start the http server
	//
	// once httpServer.Shutdown is called, it will always return a
	// http.ErrServerClosed error and we don't care about that error.
	logger.InfoContext(ctx, "listening", slog.String("address", httpServer.Addr))
	if err = httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("[in main.run] failed to listen and serve: %w", err)
	}

	// block until the server is shut down or an error occurs
	select {
	case err = <-errChan:
		return err
	case <-ctx.Done():
		return nil
	}
}
