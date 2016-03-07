
// Xavier provides a simple tool to monitor various services.
package xavier

import (
	"sync"
    "time"
    "net/http"
    "github.com/influxdb/influxdb/client/v2"
    "log"
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

// Conf type is the single source of truth for configuration of
// the system.
type Conf struct {
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

func main() {
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

	var serviceListConf = map[string]ServiceConf{
		"Github": ServiceConf{url: "http://github.com"},
		"Rbox":   ServiceConf{url: "http://www.recruiterbox.com"},
		"Google": ServiceConf{url: "http://www.google.com"},
		"Reddit": ServiceConf{url: "http://www.reddit.com"},
	}

	testConf := &Conf{
		serviceList: serviceListConf,
		timeout:     time.Second * 10,
	}

	wg.Add(1)

	go MonitorService(testConf, batchPoints, influxClient, &wg)

	wg.Wait()
}
