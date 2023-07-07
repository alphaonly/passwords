package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	conf "passwords/internal/configuration/client"
	"passwords/internal/pkg/client"
	grpcclient "passwords/internal/pkg/client/grpc/client"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	//Build tags
	// common.PrintBuildTags(buildVersion, buildDate, buildCommit)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//Configuration parameters
	cfg := conf.NewClientConf(conf.UpdateCCFromEnvironment, conf.UpdateCCFromFlags)

	//grpc client
	grpcClient := grpcclient.NewGRPCClient(cfg.Address)

	//load public key pem file
	// var cm cryptoCommon.AgentCertificateManager
	// if ac.CryptoKey != "" {
	// 	buf, err := crypto.ReadPublicKeyFile(ac)
	// 	logging.LogFatal(err)
	// 	//get public key for cm
	// 	cm = crypto.NewRSA().ReceivePublic(buf)
	// 	logging.LogFatal(cm.Error())
	// }
	//Run client
	client.NewClient(cfg, grpcClient).Run(ctx)
	//wait SIGKILL
	channel := make(chan os.Signal, 1)
	//Graceful shutdown
	signal.Notify(channel, os.Interrupt)

	select {
	case <-channel:
		log.Print("Agent shutdown by os signal")
	case <-ctx.Done():
		log.Print("Agent shutdown by cancelled context")
	}
}
