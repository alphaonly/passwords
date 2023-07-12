package client

import (
	"context"
	"fmt"
	"os"
	"passwords/internal/client/service"
	configuration "passwords/internal/configuration/client"
	"strings"
	"time"
)

func CovertStatusError(se error) error {
	for _, v := range service.Errors {
		if strings.Contains(se.Error(), v.Error()) {
			return fmt.Errorf("status converted error: %w", v)
		}
	}
	return nil
}

// CheckAuthorization - permanently check user authorization
func CheckAuthorization(ctx context.Context, service service.Service, cfg *configuration.ClientConfiguration) {
	//Check every interval whether client is authorized in server
	ticker := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-ticker.C:
			{
				//check user is authorized
				err := service.CheckUserAuthorization(ctx, cfg.Login)
				if err != nil {
					//if user was kick out the server, cancel context
					ClearScreen()
					fmt.Printf("User %v was disconnected of server by timeout", cfg.Login)
					os.Exit(1)

				}
			}
		case <-ctx.Done():
			return
		}
	}

}
