package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/influxdb/influxdb/client/v2"
)

// MonitorService which calls endpoints constantly
func MonitorService(conf *Conf, batchPoints client.BatchPoints, influxClient client.Client, wg *sync.WaitGroup) {
	requestTimeout := conf.Timeout

	go func(conf *Conf, timeout time.Duration) {
		for {
			log.Println("Starting next batch of requests")
			for label, config := range conf.Services {

				log.Println("Making request to: ", label)

				startTime := time.Now()

				resp, err := requestTimer(&config, timeout)

				if err != nil {
					log.Fatal("Something bad happened", err)
					panic(err)
				}
				endTime := time.Now()
				totalTime := endTime.Sub(startTime)
				persisData(resp, totalTime, label, batchPoints)
			}

			log.Println("Writing to DB")
			influxClient.Write(batchPoints)
			time.Sleep(10 * time.Second)
		}

		wg.Done()
	}(conf, requestTimeout)
}

// requestTimer takes a URL builds the request and returns the result
func requestTimer(service *Service, timeout time.Duration) (resp *http.Response, err error) {
	Client := &http.Client{
		Timeout: timeout,
	}

	url := service.URL
	request, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Fatalln("Something bad happened: ", err)
	}

	request.Header.Set("User-Agent", "Xavier monitoring spider v0.1.1")

	resp, err = Client.Do(request)
	return resp, err
}

// persisData is a helper function to persist the generated response.
func persisData(resp *http.Response, execTime time.Duration, label string, batchPoints client.BatchPoints) {
	tags := map[string]string{"service": label}

	fields := map[string]interface{}{
		"latency": execTime,
		"status":  resp.Status,
	}

	point, err := client.NewPoint("serviceMonitor", tags, fields)

	if err != nil {
		log.Println("Error: ", err)
	}

	batchPoints.AddPoint(point)
}
