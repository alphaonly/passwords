// main - package main function of client
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	client "passwords/internal/client"
	grpcclient "passwords/internal/client/grpc/client"
	menu "passwords/internal/client/menu"
	clientService "passwords/internal/client/service"
	"passwords/internal/common"
	conf "passwords/internal/configuration/client"
	userDomain "passwords/internal/domain/user"
	"passwords/internal/pkg/common/logging"
	"time"
)

// tags
var buildVersion = "N/A"
var buildDate = "N/A"

func main() {
	//Build tags
	common.PrintBuildTags(buildVersion, buildDate)
	//Initiate
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//Configuration parameters
	cfg := conf.NewClientConf(conf.UpdateCCFromEnvironment, conf.UpdateCCFromFlags)
	//grpc client deals with server
	grpcClient := grpcclient.NewGRPCClient(cfg.Address)
	//Client service requests server
	service := client.NewClientService(cfg, grpcClient)

	//Work
	//Authorize user by given parameters
	user, err := service.AuthorizeUser(ctx, cfg.Login, cfg.Password)
	if err != nil {
		switch true {
		case errors.Is(err, clientService.ErrNoUserOrPassword):
			log.Fatal("No user and password ")
		case errors.Is(err, clientService.ErrUserAlreadyAuthorized):
			log.Fatalf("user %v already authorized, wait for while and try again", cfg.Login)
		case errors.Is(err, clientService.ErrWrongUserOrPassword):
			log.Fatal("wrong user or password, try again")
		case errors.Is(err, clientService.ErrNoUserExists):
			//create new user
			{
				u := userDomain.User{User: cfg.Login, Password: cfg.Password}
				err = service.AddNewUser(ctx, u)
				if err != nil {
					logging.LogFatal(err)
				}
				user = &u
				log.Printf("there was no a such user, so the user %v has been created", user.User)
				time.Sleep(3 * time.Second)
			}
		}
	}
	//Permanently checks user authorization
	go client.CheckAuthorization(ctx, service, cfg)

	//Start menu
	err = menu.New(user, cfg, service, nil).Start(ctx)
	logging.LogFatal(err)
	//End user authorization in server
	err = service.DisAuthorizeUser(ctx, cfg.Login, cfg.Password)
	logging.LogFatal(err)
	fmt.Printf("User %v disauthorized from server\n", user.User)

	fmt.Println("Exit passwords")
}
