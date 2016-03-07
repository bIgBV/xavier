package xavier

import (
    "log"
    "net/http"
    "time"	
    "github.com/influxdb/influxdb/client/v2"
    "sync"
)

// MonitorService which calls endpoints constantly
func MonitorService(conf *Conf, batchPoints client.BatchPoints, influxClient client.Client, wg *sync.WaitGroup) {
	go func(conf *Conf) {
		for {
			log.Println("Starting next batch of requests")
			for label, config := range conf.serviceList {

				url := config.url
				
				log.Println("Making request")

				startTime := time.Now()
				resp, err := requestTimer(url)
                if err != nil {
                    log.Fatal(err)
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
	}(conf)
}

func requestTimer(url string) (resp *http.Response, err error) {
    Client := &http.Client{}
    request, err := http.NewRequest("HEAD", url, nil)
    if err != nil {
        log.Fatalln(err)
    }

    request.Header.Set("User-Agent", "Xavier monitoring spider v0.1")
    
    resp, err = Client.Do(request)
    return resp, err
}

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