package app

import (
	"log/slog"
	"time"

	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/services/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	//init storage
	stotage, err := sqlite.New(storagePath)
	if err != nil {
		log.Error("failed to init storage", slog.String("path", storagePath), err)
		panic(err)
	}
	//init services
	authService := auth.New(log, stotage, stotage, stotage, tokenTTL)	
	//TODO init auth server
	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{
		GRPCSrv: grpcApp,
	}
}
