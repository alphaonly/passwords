package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"


	"strconv"
	"time"

	"math/rand"
	"net/url"

	conf "passwords/internal/configuration/client"
)

t
type Client struct {
	Configuration *conf.AgentConfiguration
	GRPCClient    *client.GRPCClient
}

func NewClient(c *conf.AgentConfiguration,grpcClient *client.GRPCClient) Client {

	return Client{
		Configuration: c,
		GRPCClient:   grpcClient,
	}
}


// SendDataGRPC - sends batch metric data using gRPC client in stream
func SendDataGRPC(ctx context.Context, grpcClient *client.GRPCClient) error {


	var wg sync.WaitGroup
	//get stream
	stream, err := grpcClient.Client.AddMetricMulti(ctx)
	logging.LogFatal(err)
	//iterate the array on every metric data
	for _, metric := range *sd.JSONBatchBody {
		//Send data in parallel
		wg.Add(1)
		go func(metric schema.Metrics) {
			//make metric gRPC structure
			protoMetric := &proto.Metric{
				Name: metric.ID,
				Type: common.ConvertMetricType(metric.MType),
			}
			//determine which one of metric param is fulfilled
			switch {
			case metric.Value != nil:
				protoMetric.Gauge = *metric.Value
			case metric.Delta != nil:
				protoMetric.Counter = *metric.Delta
			}
			//send data
			err = stream.Send(&proto.AddMetricRequest{Metric: protoMetric})
			logging.LogPrintln(err)
			//mark routine as finished
			wg.Done()
		}(metric)
	}
	//wait for every routine is finished
	wg.Wait()

	//Capture response
	var resp *proto.AddMetricResponse
	go func(resp *proto.AddMetricResponse) {

		wg.Add(1)
		for {
			//receive response
			resp, err = stream.Recv()
			if err == io.EOF {
				log.Println("getting response is finished, everything is well")
				wg.Done()
				break
			}
			logging.LogFatal(err)
		}
	}(resp)

	wg.Wait()

	return err
}
func (a Client) Update(ctx context.Context, metrics *Metrics) {
	var m runtime.MemStats

	ticker := time.NewTicker(time.Duration(a.Configuration.PollInterval))

	defer ticker.Stop()
repeatAgain:
	select {
	case <-ticker.C:
		{
			runtime.ReadMemStats(&m)

			a.UpdateLocker.Lock()

			metrics.Alloc = Gauge(m.Alloc)
			metrics.BuckHashSys = Gauge(m.BuckHashSys)
			metrics.Frees = Gauge(m.Frees)
			metrics.GCCPUFraction = Gauge(m.GCCPUFraction)
			metrics.GCSys = Gauge(m.GCSys)
			metrics.HeapAlloc = Gauge(m.HeapAlloc)
			metrics.HeapIdle = Gauge(m.HeapIdle)
			metrics.HeapInuse = Gauge(m.HeapInuse)
			metrics.HeapObjects = Gauge(m.HeapObjects)
			metrics.HeapReleased = Gauge(m.HeapReleased)
			metrics.HeapSys = Gauge(m.HeapSys)
			metrics.LastGC = Gauge(m.LastGC)
			metrics.Lookups = Gauge(m.Lookups)
			metrics.MCacheInuse = Gauge(m.MCacheInuse)
			metrics.MCacheSys = Gauge(m.MCacheSys)
			metrics.MSpanInuse = Gauge(m.MSpanInuse)
			metrics.MSpanSys = Gauge(m.MSpanSys)
			metrics.Mallocs = Gauge(m.Mallocs)
			metrics.NextGC = Gauge(m.NextGC)
			metrics.NumForcedGC = Gauge(m.NumForcedGC)
			metrics.NumGC = Gauge(m.NumGC)
			metrics.OtherSys = Gauge(m.OtherSys)
			metrics.PauseTotalNs = Gauge(m.PauseTotalNs)
			metrics.StackInuse = Gauge(m.StackInuse)
			metrics.StackSys = Gauge(m.StackSys)
			metrics.Sys = Gauge(m.Sys)
			metrics.TotalAlloc = Gauge(m.TotalAlloc)
			metrics.RandomValue = Gauge(rand.Int63())
			metrics.PollCount++

			a.UpdateLocker.Unlock()
			goto repeatAgain
		}
	case <-ctx.Done():
		{
			log.Println("Metrics reading cancelled by context")
			return
		}
	}

}

func (a Client) UpdateGOPS(ctx context.Context, metrics *Metrics) {

	ticker := time.NewTicker(time.Duration(a.Configuration.PollInterval))
	defer ticker.Stop()
repeatAgain:
	select {
	case <-ticker.C:
		{
			v, err := mem.VirtualMemory()
			if err != nil {
				log.Println(err)
				ctx.Done()
			}
			i, err := cpu.Times(true)
			if err != nil {
				log.Println(err)
				ctx.Done()
			}
			a.UpdateLocker.Lock()
			metrics.TotalMemory = Gauge(v.Total)
			metrics.FreeMemory = Gauge(v.Free)
			metrics.CPUutilization1 = Gauge(i[0].System)
			a.UpdateLocker.Unlock()

			goto repeatAgain
		}
	case <-ctx.Done():
		{
			log.Println("context is done received, metrics reading cancelled by app context")
			return
		}
	}

}

func (a Client) CompressData(data map[*sender]bool) map[*sender]bool {

	switch a.Configuration.CompressType {

	case "gzip":
		{
			var body any
			for k := range data {
				if k.JSONBody != nil {
					body = *k.JSONBody
				} else if k.JSONBatchBody != nil {
					body = *k.JSONBatchBody
				} else {
					logFatal(errors.New("agent:nothing to marshal as sendData bodies are nil"))
				}

				b, err := json.Marshal(body)
				logFatal(err)

				k.compressedBody, err = compression.GzipCompress(b)
				logFatal(err)
			}
		}
	}

	return data
}
func (a Client) prepareData(metrics *Metrics) map[*sender]bool {
	m := make(map[*sender]bool)
	keys := make(HeaderKeys)

	keys["X-Real-IP"] = a.Configuration.Address

	switch a.Configuration.CompressType {
	case "deflate":
		{
			keys["Accept-Encoding"] = "deflate"
			keys["Content-Encoding"] = "deflate"
		}
	case "gzip":
		{
			keys["Accept-Encoding"] = "gzip"
			keys["Content-Encoding"] = "gzip"
		}
	}

	switch a.Configuration.Mode {

	case conf.MODE_JSON_BATCH, conf.MODE_GRPC: //JSON Batch or gRPC
		{
			keys["Content-Type"] = "application/json"
			keys["Accept"] = "application/json"

			data := &sender{
				url:    a.baseURL.JoinPath("updates"),
				keys:   keys,
				signer: a.Signer,
			}

			AddGaugeDataJSONToBatch(data, metrics.Alloc, "Alloc")
			AddGaugeDataJSONToBatch(data, metrics.GCCPUFraction, "GCCPUFraction")
			AddGaugeDataJSONToBatch(data, metrics.GCSys, "GCSys")
			AddGaugeDataJSONToBatch(data, metrics.HeapAlloc, "HeapAlloc")
			AddGaugeDataJSONToBatch(data, metrics.HeapIdle, "HeapIdle")
			AddGaugeDataJSONToBatch(data, metrics.HeapInuse, "HeapInuse")
			AddGaugeDataJSONToBatch(data, metrics.HeapObjects, "HeapObjects")
			AddGaugeDataJSONToBatch(data, metrics.HeapReleased, "HeapReleased")
			AddGaugeDataJSONToBatch(data, metrics.HeapSys, "HeapSys")
			AddGaugeDataJSONToBatch(data, metrics.LastGC, "LastGC")
			AddGaugeDataJSONToBatch(data, metrics.Lookups, "Lookups")
			AddGaugeDataJSONToBatch(data, metrics.MCacheSys, "MCacheSys")
			AddGaugeDataJSONToBatch(data, metrics.MSpanInuse, "MSpanInuse")
			AddGaugeDataJSONToBatch(data, metrics.MSpanSys, "MSpanSys")
			AddGaugeDataJSONToBatch(data, metrics.Mallocs, "Mallocs")
			AddGaugeDataJSONToBatch(data, metrics.NextGC, "NextGC")
			AddGaugeDataJSONToBatch(data, metrics.NumForcedGC, "NumForcedGC")
			AddGaugeDataJSONToBatch(data, metrics.NumGC, "NumGC")
			AddGaugeDataJSONToBatch(data, metrics.OtherSys, "OtherSys")
			AddGaugeDataJSONToBatch(data, metrics.PauseTotalNs, "PauseTotalNs")
			AddGaugeDataJSONToBatch(data, metrics.StackInuse, "StackInuse")
			AddGaugeDataJSONToBatch(data, metrics.StackSys, "StackSys")
			AddGaugeDataJSONToBatch(data, metrics.Sys, "Sys")
			AddGaugeDataJSONToBatch(data, metrics.TotalAlloc, "TotalAlloc")
			AddGaugeDataJSONToBatch(data, metrics.RandomValue, "RandomValue")
			AddGaugeDataJSONToBatch(data, metrics.Frees, "Frees")
			AddCounterDataJSONToBatch(data, metrics.PollCount, "PollCount")

			AddGaugeDataJSONToBatch(data, metrics.TotalMemory, "TotalMemory")
			AddGaugeDataJSONToBatch(data, metrics.FreeMemory, "FreeMemory")
			AddGaugeDataJSONToBatch(data, metrics.CPUutilization1, "CPUutilization1")

			//Encrypt data to send with TLS public key
			if a.Configuration.CryptoKey != "" {
				bts, err := json.Marshal(&data.JSONBatchBody)
				logging.LogFatal(err)
				data.encryptedBody = a.Cm.EncryptData(bts)
				logging.LogFatal(a.Cm.Error())
			}

			m[data] = true
		}
	case conf.MODE_JSON: //JSON
		{
			keys["Content-Type"] = "application/json"
			keys["Accept"] = "application/json"

			data := sender{
				url:    a.baseURL.JoinPath("update"),
				keys:   keys,
				signer: a.Signer,
			}
			AddGaugeDataJSON(data, metrics.Alloc, "Alloc", m)
			AddGaugeDataJSON(data, metrics.Frees, "Frees", m)
			AddGaugeDataJSON(data, metrics.GCCPUFraction, "GCCPUFraction", m)
			AddGaugeDataJSON(data, metrics.GCSys, "GCSys", m)
			AddGaugeDataJSON(data, metrics.HeapAlloc, "HeapAlloc", m)
			AddGaugeDataJSON(data, metrics.HeapIdle, "HeapIdle", m)
			AddGaugeDataJSON(data, metrics.HeapInuse, "HeapInuse", m)
			AddGaugeDataJSON(data, metrics.HeapObjects, "HeapObjects", m)
			AddGaugeDataJSON(data, metrics.HeapReleased, "HeapReleased", m)
			AddGaugeDataJSON(data, metrics.HeapSys, "HeapSys", m)
			AddGaugeDataJSON(data, metrics.LastGC, "LastGC", m)
			AddGaugeDataJSON(data, metrics.Lookups, "Lookups", m)
			AddGaugeDataJSON(data, metrics.MCacheSys, "MCacheSys", m)
			AddGaugeDataJSON(data, metrics.MSpanInuse, "MSpanInuse", m)
			AddGaugeDataJSON(data, metrics.MSpanSys, "MSpanSys", m)
			AddGaugeDataJSON(data, metrics.Mallocs, "Mallocs", m)
			AddGaugeDataJSON(data, metrics.NextGC, "NextGC", m)
			AddGaugeDataJSON(data, metrics.NumForcedGC, "NumForcedGC", m)
			AddGaugeDataJSON(data, metrics.NumGC, "NumGC", m)
			AddGaugeDataJSON(data, metrics.OtherSys, "OtherSys", m)
			AddGaugeDataJSON(data, metrics.PauseTotalNs, "PauseTotalNs", m)
			AddGaugeDataJSON(data, metrics.StackInuse, "StackInuse", m)
			AddGaugeDataJSON(data, metrics.StackSys, "StackSys", m)
			AddGaugeDataJSON(data, metrics.Sys, "Sys", m)
			AddGaugeDataJSON(data, metrics.TotalAlloc, "TotalAlloc", m)
			AddGaugeDataJSON(data, metrics.RandomValue, "RandomValue", m)
			AddCounterDataJSON(data, metrics.PollCount, "PollCount", m)

			//GOPSUtil
			AddGaugeDataJSON(data, metrics.TotalMemory, "TotalMemory", m)
			AddGaugeDataJSON(data, metrics.FreeMemory, "FreeMemory", m)
			AddGaugeDataJSON(data, metrics.CPUutilization1, "CPUutilization1", m)

			// Encrypt data to send with TLS public key
			if a.Configuration.CryptoKey != "" {
				bts, err := json.Marshal(&data.JSONBody)
				logging.LogFatal(err)
				data.encryptedBody = a.Cm.EncryptData(bts)
				logging.LogFatal(a.Cm.Error())
			}
		}
	case conf.MODE_NO_JSON:
		{

			keys["Content-Type"] = "plain/text"
			keys["Accept"] = "text/html"

			data := sender{
				url:  a.baseURL.JoinPath("update"),
				keys: keys,
			}

			AddGaugeData(data, metrics.Alloc, "Alloc", m)
			AddGaugeData(data, metrics.Frees, "Frees", m)
			AddGaugeData(data, metrics.GCCPUFraction, "GCCPUFraction", m)
			AddGaugeData(data, metrics.GCSys, "GCSys", m)
			AddGaugeData(data, metrics.HeapAlloc, "HeapAlloc", m)
			AddGaugeData(data, metrics.HeapIdle, "HeapIdle", m)
			AddGaugeData(data, metrics.HeapInuse, "HeapInuse", m)
			AddGaugeData(data, metrics.HeapObjects, "HeapObjects", m)
			AddGaugeData(data, metrics.HeapReleased, "HeapReleased", m)
			AddGaugeData(data, metrics.HeapSys, "HeapSys", m)
			AddGaugeData(data, metrics.LastGC, "LastGC", m)
			AddGaugeData(data, metrics.Lookups, "Lookups", m)
			AddGaugeData(data, metrics.MCacheSys, "MCacheSys", m)
			AddGaugeData(data, metrics.MSpanInuse, "MSpanInuse", m)
			AddGaugeData(data, metrics.MSpanSys, "MSpanSys", m)
			AddGaugeData(data, metrics.Mallocs, "Mallocs", m)
			AddGaugeData(data, metrics.NextGC, "NextGC", m)
			AddGaugeData(data, metrics.NumForcedGC, "NumForcedGC", m)
			AddGaugeData(data, metrics.NumGC, "NumGC", m)
			AddGaugeData(data, metrics.OtherSys, "OtherSys", m)
			AddGaugeData(data, metrics.PauseTotalNs, "PauseTotalNs", m)
			AddGaugeData(data, metrics.StackInuse, "StackInuse", m)
			AddGaugeData(data, metrics.StackSys, "StackSys", m)
			AddGaugeData(data, metrics.Sys, "Sys", m)
			AddGaugeData(data, metrics.TotalAlloc, "TotalAlloc", m)
			AddGaugeData(data, metrics.RandomValue, "RandomValue", m)
			AddCounterData(data, metrics.PollCount, "PollCount", m)
			//GOPSUtil
			AddGaugeData(data, metrics.TotalMemory, "TotalMemory", m)
			AddGaugeData(data, metrics.FreeMemory, "FreeMemory", m)
			AddGaugeData(data, metrics.CPUutilization1, "CPUutilization1", m)
		}

	}

	return m
}

func (a Client) Send(ctx context.Context, metrics *Metrics) {

	ticker := time.NewTicker(time.Duration(a.Configuration.ReportInterval))

	defer ticker.Stop()

	//Worker pool
	var f workerpool.TypicalJobFunction[*sender]
	var wp workerpool.WorkerPool[*sender]

	if a.Configuration.RateLimit > 0 {
		//Initialize a job function for workerpool
		f = func(key *sender) workerpool.JobResult {
			err := key.SendDataResty(a.RestyClient)
			if err != nil {
				log.Println(err)
				return workerpool.JobResult{Result: err.Error()}
			}

			return workerpool.JobResult{Result: "OK"}
		}
		wp = workerpool.NewWorkerPool[*sender](a.Configuration.RateLimit)
		// worker pool start
		wp.Start(ctx)
	}
repeatAgain:
	select {
	case <-ticker.C:
		{
			dataPackage := a.prepareData(metrics)
			dataPackage = a.CompressData(dataPackage)
			i := 0
			for key := range dataPackage {
				i++
				name := fmt.Sprintf("Send metric job %v", i)
				switch a.Configuration.RateLimit {
				case 0, 1:
					{
						//Client depends on Mode
						switch a.Configuration.Mode {
						case conf.MODE_GRPC: //gRPC Mode
							{
								err := key.SendDataGRPC(ctx, a.GRPCClient)
								if err != nil {
									log.Println(err)
									return
								}

							}
						default: //No gRPC MOde
							{
								err := key.SendDataResty(a.RestyClient)
								if err != nil {
									log.Println(err)
									return
								}

							}
						}
					}
				default:
					{
						//send a job to worker pool
						job := workerpool.Job[*sender]{Name: name, Data: key, Func: f}
						wp.SendJob(ctx, job)
					}
				}
			}
			goto repeatAgain
		}
	case <-ctx.Done():
		{

			a.GRPCClient.Close()
			break
		}
	}

}

func (a Client) Run(ctx context.Context) {

	metrics := Metrics{}
	go a.Update(ctx, &metrics)
	go a.UpdateGOPS(ctx, &metrics)
	go a.Send(ctx, &metrics)

}
