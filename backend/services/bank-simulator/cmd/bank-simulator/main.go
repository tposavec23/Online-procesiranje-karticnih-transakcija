// Command bank-simulator simulates Bank 1: an H2H (ISO 8583) server that approves every sale.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sevenpay/backend/services/bank-simulator/internal/simulator"
)

func main() {
	if err := run(context.Background(), os.Getenv, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "bank-simulator: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, getenv func(string) string, stdout io.Writer) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	port := getenv("PORT")
	if port == "" {
		return errors.New("PORT is not set")
	}
	logger := slog.New(slog.NewJSONHandler(stdout, nil))

	srv, err := simulator.NewH2HServer(logger)
	if err != nil {
		return err
	}
	if err := srv.Start(":" + port); err != nil {
		return fmt.Errorf("listen on port %s: %w", port, err)
	}
	logger.Info("bank simulator listening", "addr", srv.Addr)

	<-ctx.Done()
	stop()
	logger.Info("bank simulator shutting down")
	srv.Close()

	return nil
}
