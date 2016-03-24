package main

import (
	"log"

	"github.com/influxdb/influxdb/client/v2"
)

type Persist interface {
	PersistData(outBuf chan<- *ServiceResp)
}

// persisData is a helper function to persist the generated response.
func persisData(resp *ServiceResp, batchPoints client.BatchPoints) {
	tags := map[string]string{"service": resp.label}

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
