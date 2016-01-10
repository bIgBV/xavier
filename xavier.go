// Package Xavier provides a simple tool to monitor various services.
package main

import (
	"github.com/influxdb/influxdb/client/v2"
	"log"
	"net/http"
	"time"
)

// XavierResponse type is what is sent to persistance adapters.
// Since the adapters can be configured to persist whatever data from
// the response they choose to, it contains a reference to the response
// for a particular request.
type XavierResponse struct {
	resp    *http.Response
	label   string
	latency time.Duration
}

// XavierConf type is the single source of truth for configuration of
// the system.
type XavierConf struct {
	serviceList map[string]ServiceConf
	timeout     time.Duration
}

// ServiceConf type stores configuration regarding indidual services to be
// monitored. This is data such as the URL, auth, timeout, etc.
type ServiceConf struct {
	url string
}

const (
	MyDB     = "MonitorData"
	username = "xavier"
	password = "watcheverything"
)

func MonitorService(conf *XavierConf, batchPoints client.BatchPoints, influxClient client.Client) {
	Client := &http.Client{}
	tempChan := make(chan string)

	go func(tempChan chan string, conf *XavierConf) {
		for {
			log.Println("Starting next batch of requests")
			for label, config := range conf.serviceList {

				url := config.url

				request, err := http.NewRequest("HEAD", url, nil)
				if err != nil {
					log.Fatalln(err)
				}

				request.Header.Set("User-Agent", "Xavier monitoring spider v0.1")

				log.Println("Making request")

				startTime := time.Now()
				response, err := Client.Do(request)
				if err != nil {
					log.Fatalln(err)
				}
				endTime := time.Now()
				totalTime := endTime.Sub(startTime)

				tags := map[string]string{"service": label}
				fields := map[string]interface{}{
					"latency": totalTime,
					"status":  response.Status,
				}

				point, err := client.NewPoint("serviceMonitor", tags, fields)

				if err != nil {
					log.Println("Error: ", err)
				}

				tempChan <- "ping"

				batchPoints.AddPoint(point)

			}
			log.Println("Writing to DB")
			influxClient.Write(batchPoints)
			time.Sleep(10 * time.Second)

		}
	}(tempChan, conf)
	for msg := range tempChan {
		log.Println(msg)
	}
}

func main() {

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Println("Error: ", err)
	}

	batchPoints, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})

	if err != nil {
		log.Println("Error: ", err)
	}

	var serviceListConf = map[string]ServiceConf{
		"Github": ServiceConf{url: "http://github.com"},
		"Rbox":   ServiceConf{url: "http://www.recruiterbox.com"},
		"Google": ServiceConf{url: "http://www.google.com"},
		"Reddit": ServiceConf{url: "http://www.reddit.com"},
	}

	testConf := &XavierConf{
		serviceList: serviceListConf,
		timeout:     time.Second * 10,
	}

	MonitorService(testConf, batchPoints, influxClient)
}
