package client

import (
	"context"
	"log"

	proto "passwords/internal/adapters/grpc/proto"
	conf "passwords/internal/configuration/client"
	accountDomain "passwords/internal/domain/account"
	userDomain "passwords/internal/domain/user"
	grpcCl "passwords/internal/pkg/client/grpc/client"

	service "passwords/internal/pkg/client/service"
)

type myClient struct {
	Configuration *conf.ClientConfiguration
	GRPCClient    *grpcCl.GRPCClient
}

func NewClient(c *conf.ClientConfiguration, grpcClient *grpcCl.GRPCClient) service.Service {

	return &myClient{
		Configuration: c,
		GRPCClient:    grpcClient,
	}
}

func (c myClient) AuthorizeUser(ctx context.Context, user string, password string) error {

	response, err := c.GRPCClient.Client.GetUser(ctx, &proto.GetUserRequest{Login: user})
	if err != nil {
		return service.ErrWrongUserOrPassword
	}
	log.Printf("User %v authorized ",response.User.Name)

	return nil
}
func (c myClient) AddNewUser(ctx context.Context, user userDomain.User) error             { return nil }
func (c myClient) AddNewAccount(ctx context.Context, account accountDomain.Account) error { return nil }
func (c myClient) GetAllAccounts(ctx context.Context, user string) (accountDomain.Accounts, error) {
	return nil, nil
}

// // SendDataGRPC - sends batch metric data using gRPC client in stream
// func SendDataGRPC(ctx context.Context, grpcClient *client.GRPCClient) error {
// 	var wg sync.WaitGroup
// 	//get stream
// 	stream, err := grpcClient.Client.AddMetricMulti(ctx)
// 	logging.LogFatal(err)
// 	//iterate the array on every metric data
// 	for _, metric := range *sd.JSONBatchBody {
// 		//Send data in parallel
// 		wg.Add(1)
// 		go func(metric schema.Metrics) {
// 			//make metric gRPC structure
// 			protoMetric := &proto.Metric{
// 				Name: metric.ID,
// 				Type: common.ConvertMetricType(metric.MType),
// 			}
// 			//determine which one of metric param is fulfilled
// 			switch {
// 			case metric.Value != nil:
// 				protoMetric.Gauge = *metric.Value
// 			case metric.Delta != nil:
// 				protoMetric.Counter = *metric.Delta
// 			}
// 			//send data
// 			err = stream.Send(&proto.AddMetricRequest{Metric: protoMetric})
// 			logging.LogPrintln(err)
// 			//mark routine as finished
// 			wg.Done()
// 		}(metric)
// 	}
// 	//wait for every routine is finished
// 	wg.Wait()

// 	//Capture response
// 	var resp *proto.AddMetricResponse
// 	go func(resp *proto.AddMetricResponse) {

// 		wg.Add(1)
// 		for {
// 			//receive response
// 			resp, err = stream.Recv()
// 			if err == io.EOF {
// 				log.Println("getting response is finished, everything is well")
// 				wg.Done()
// 				break
// 			}
// 			logging.LogFatal(err)
// 		}
// 	}(resp)

// 	wg.Wait()

// 	return err
// }

func (c myClient) Run(ctx context.Context) {

}

// func (a Client) CompressData(data map[*sender]bool) map[*sender]bool {

// 	switch a.Configuration.CompressType {

// 	case "gzip":
// 		{
// 			var body any
// 			for k := range data {
// 				if k.JSONBody != nil {
// 					body = *k.JSONBody
// 				} else if k.JSONBatchBody != nil {
// 					body = *k.JSONBatchBody
// 				} else {
// 					logFatal(errors.New("agent:nothing to marshal as sendData bodies are nil"))
// 				}

// 				b, err := json.Marshal(body)
// 				logFatal(err)

// 				k.compressedBody, err = compression.GzipCompress(b)
// 				logFatal(err)
// 			}
// 		}
// 	}

// 	return data
// }

func (c myClient) Stop() {

}
