package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/influxdb/influxdb/client/v2"
)

// ServiceResp is a datatype for storing response data to be persisted
type ServiceResp struct {
	Resp      *http.Response
	Latency   time.Duration
	label     string
	TimeStamp time.Time
}

// Label is a simple getter func for the timestamp of the response object
func (resp *ServiceResp) Label() string {
	return resp.label
}

// SetLabel is a setter for timestamp for the response object
func (resp *ServiceResp) SetLabel(l string) {
	resp.label = l
}

// MonitorService which calls endpoints constantly
func MonitorService(conf *Conf, batchPoints client.BatchPoints, influxClient client.Client, wg *sync.WaitGroup) {

	go func(conf *Conf) {
		for {
			log.Println("Starting next batch of requests")
			for label, config := range conf.Services {

				log.Println("Making request to: ", label)

				resp := timeRequest(&config)

				resp.SetLabel(label)

				persisData(&resp, batchPoints)
			}

			log.Println("Writing to DB")
			influxClient.Write(batchPoints)
			time.Sleep(5 * time.Second)
		}

		wg.Done()
	}(conf)
}

// timeRequest takes a URL builds the request and returns the result
func timeRequest(service *Service) (serviceResp ServiceResp) {
	Client := &http.Client{
		Timeout: service.Timeout,
	}

	startTime := time.Now()

	url := service.URL
	request, err := http.NewRequest("HEAD", url, nil)

	if err != nil {
		log.Fatalln("Something bad happened: ", err)
	}

	request.Header.Set("User-Agent", "Xavier monitoring spider v0.1.1")

	resp, err := Client.Do(request)

	if err != nil {
		log.Fatal("Something bad happened", err)
		panic(err)
	}

	endTime := time.Now()
	totalTime := endTime.Sub(startTime)

	serviceResp = ServiceResp{
		Resp:      resp,
		Latency:   totalTime,
		TimeStamp: time.Now(),
	}

	return serviceResp
}
