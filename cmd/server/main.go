package main

import (
	"context"
	"os"
	"os/signal"

	grpcService "passwords/internal/adapters/grpc"
	"passwords/internal/composites"
	configuration "passwords/internal/configuration/server"
	"passwords/internal/pkg/dbclient/postgres"
	"passwords/internal/pkg/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//configuration
	cfg := configuration.NewServerConf(configuration.UpdateSCFromEnvironment, configuration.UpdateSCFromFlags)

	//postgres storage
	DBclient := postgres.NewPostgresClient(ctx, cfg.DatabaseURI)

	//composites
	userComposite := composites.NewUserComposite(DBclient, cfg)
	accountComposite := composites.NewAccountComposite(DBclient, userComposite.Service, cfg)

	//grpc
	grpcService := grpcService.NewGRPCService(userComposite.Storage, accountComposite.Storage)
	
	//server
	srv := server.NewServer(cfg, grpcService)

	go srv.Run()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	<-osSignal
	srv.Stop()

}
