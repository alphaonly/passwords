package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"passwords/internal/client/crypto"
	"passwords/internal/common"

	// cryptoCommon "passwords/internal/common/crypto"
	"passwords/internal/pkg/common/logging"

	grpcclient "passwords/internal/agent/grpc/client"
	conf "passwords/internal/configuration"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	//Build tags
	common.PrintBuildTags(buildVersion, buildDate, buildCommit)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//Configuration parameters
	ac := conf.NewAgentConf(conf.UpdateACFromEnvironment, conf.UpdateACFromFlags)

	//grpc client
	grpcClient := grpcclient.NewGRPCClient(ac.Address)

	//load public key pem file
	// var cm cryptoCommon.AgentCertificateManager
	// if ac.CryptoKey != "" {
	// 	buf, err := crypto.ReadPublicKeyFile(ac)
	// 	logging.LogFatal(err)
	// 	//get public key for cm
	// 	cm = crypto.NewRSA().ReceivePublic(buf)
	// 	logging.LogFatal(cm.Error())
	// }
	//Run agent
	client.NewAgent(ac, grpcClient, cm).Run(ctx)
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
