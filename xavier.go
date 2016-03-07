// Package xavier provides a simple tool to monitor various services.
package xavier

import (
	"github.com/influxdb/influxdb/client/v2"
	"log"
	"net/http"
	"sync"
	"time"
)

// Response type is what is sent to persistance adapters.
// Since the adapters can be configured to persist whatever data from
// the response they choose to, it contains a reference to the response
// for a particular request.
type Response struct {
	resp    *http.Response
	label   string
	latency time.Duration
}

const (
	MyDB     = "MonitorData"
	username = "xavier"
	password = "watcheverything"
	confName = "config.yml"
)

func main() {

	log.Println("Test")
	var wg sync.WaitGroup

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

	// var serviceConfList = map[string]Service{
	// 	"Github": Service{url: "http://github.com"},
	// 	"Rbox":   Service{url: "http://www.recruiterbox.com"},
	// 	"Google": Service{url: "http://www.google.com"},
	// 	"Reddit": Service{url: "http://www.reddit.com"},
	// }

	// testConf := &Conf{
	// 	serviceList: serviceConfList,
	// 	timeout:     time.Second * 10,
	// }

	conf := parseConfig(confName)

	log.Println(conf)

	wg.Add(1)

	go MonitorService(conf, batchPoints, influxClient, &wg)

	wg.Wait()
}
