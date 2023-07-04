package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"passwords/internal/composites"
	configuration "passwords/internal/configuration/server"
	"passwords/internal/pkg/common/logging"
	"passwords/internal/pkg/dbclient/postgres"
	"passwords/internal/pkg/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := configuration.NewServerConf(configuration.UpdateSCFromEnvironment, configuration.UpdateSCFromFlags)

	dbclient := postgres.NewPostgresClient(ctx, cfg.DatabaseURI)

	userComposite := composites.NewUserComposite(dbclient, cfg)
	accountComposite := composites.NewAccountComposite(dbclient, userComposite.Service, cfg)

	httpServer := &http.Server{
		Addr: cfg.RunAddress,
	}

	srv := server.NewServer(httpServer)

	go srv.Run()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	<-osSignal
	err := srv.Stop(ctx)
	logging.LogFatal(err)
}
