package main

import (
	"context"
	"os"
	"os/signal"
	"passwords/internal/server"
	"sync"

	grpcService "passwords/internal/adapters/grpc"
	"passwords/internal/composites"
	configuration "passwords/internal/configuration/server"
	"passwords/internal/pkg/dbclient/postgres"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//configuration
	cfg := configuration.NewServerConf(configuration.UpdateSCFromEnvironment, configuration.UpdateSCFromFlags)

	//postgres storage
	DBClient := postgres.NewPostgresClient(ctx, cfg.DatabaseURI)

	//composites
	userComposite := composites.NewUserComposite(DBClient)
	accountComposite := composites.NewAccountComposite(DBClient)

	//Authorized users
	authorizedUsers := &sync.Map{}
	//Last users operation
	lastUserOperation := &sync.Map{}
	//grpc service
	gService := grpcService.NewGRPCService(userComposite.Storage, accountComposite.Storage, authorizedUsers, lastUserOperation)
	//server
	srv := server.NewServer(cfg, gService, authorizedUsers, lastUserOperation)

	//Go-routine listens  to some port given in configuration
	go srv.RunListener()
	//Go-routine stops user authorization after some interval(given in configuration) of client operations absence,
	//if there are no operations from client authorization stops
	go srv.RunUserKicker(ctx)

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)

	<-osSignal
	srv.Stop()

}
