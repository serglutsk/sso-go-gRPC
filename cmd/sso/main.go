package main

import (
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//  config init
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting app")
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	
	go application.GRPCSrv.MustRun()
	//Grecefull shutdown
	stop := make(chan os.Signal, 1)//bufered channel signal from OS
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)//Якщо прийдуть перелічені сигнали, відправ їх у мій канал stop
	//syscall.SIGINT == Ctrl+C
	//syscall.SIGTERM - signal from Kubernetes or Docker
	sign := <-stop
	log.Info("stoping application", slog.String("signal", sign.String()))
	application.GRPCSrv.Stop()
	log.Info("Gracefully stopped")
	
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})) // os.Stdout - show log in console
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))	
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))	
	}

	return log
}

