package main

import (
	"log"

	"github.com/influxdb/influxdb/client/v2"
)

type Persistance interface {
	PersistData(outBuf chan<- *ServiceResp)
}

// persisData is a helper function to persist the generated response.
func PersisData(resp *ServiceResp, batchPoints client.BatchPoints) {
	tags := map[string]string{"service": resp.label}

	log.SetOutput(io.MultiWriter(os.Stderr, createLogFile("xavier.log")))

	fields := map[string]interface{}{
		"latency": resp.Latency,
		"status":  resp.Resp.Status,
	}

	point, err := client.NewPoint("serviceMonitor", tags, fields)

	if err != nil {
		log.Println("Error: ", err)
	}

	batchPoints.AddPoint(point)
}
