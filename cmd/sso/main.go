package main

import (
	"fmt"
	"go-sso/intenal/config"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// initializing config object
	cnf := config.MustLoad() // !!functions with Must getting panic when occurred error !

	// initializing logger
	log := setupLogger(cnf.Env) // we use slog
	fmt.Println(log)
	// todo: initizlizing application (app)

	// todo: Run grps server application
}
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	// select log level depending on environment
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
