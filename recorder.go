package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "os"
  "encoding/json"
  "github.com/samuel/go-librato/librato"
)

type Recording struct {
  Queue string
  PersistentCount   float64
  AvgIngressRate    float64
  AvgEgressRate     float64
  AvgAckIngressRate float64
  AvgAckEgressRate  float64
}

func Collect(queue string) Recording {
  client := &http.Client{}

  url := os.Getenv("RABBIT_DOMAIN") + "/api/queues/" + os.Getenv("RABBIT_VHOST") + "/" + queue
  req, err := http.NewRequest("GET", url, nil)
  req.SetBasicAuth(os.Getenv("RABBIT_USER"), os.Getenv("RABBIT_PASSWORD"))

  response, err := client.Do(req)

  if err != nil {
    panic(err)
  }

  defer response.Body.Close()
  data, err := ioutil.ReadAll(response.Body)

  if err != nil {
    panic(err)
  }

  var stats interface{}
  err = json.Unmarshal(data, &stats)

  if err != nil {
    panic(err)
  }

  root := stats.(map[string]interface{})
  bqs := root["backing_queue_status"].(map[string]interface{})

  result := Recording{Queue:             queue,
                      PersistentCount:   bqs["persistent_count"].(float64),
                      AvgIngressRate:    bqs["avg_ingress_rate"].(float64),
                      AvgEgressRate:     bqs["avg_egress_rate"].(float64),
                      AvgAckIngressRate: bqs["avg_ack_ingress_rate"].(float64),
                      AvgAckEgressRate:  bqs["avg_ack_egress_rate"].(float64)}

  return result
}

func Record(stats Recording, source string) {
  var data []interface{}=[]interface{}{
    librato.Metric{Name: "cloudamqp."+stats.Queue+".persistent_count",     Value: stats.PersistentCount},
    librato.Metric{Name: "cloudamqp."+stats.Queue+".avg_ingress_rate",     Value: stats.AvgIngressRate},
    librato.Metric{Name: "cloudamqp."+stats.Queue+".avg_egress_rate",      Value: stats.AvgEgressRate},
    librato.Metric{Name: "cloudamqp."+stats.Queue+".avg_ack_ingress_rate", Value: stats.AvgAckIngressRate},
    librato.Metric{Name: "cloudamqp."+stats.Queue+".avg_ack_egress_rate",  Value: stats.AvgAckEgressRate},
  }

  client   := &librato.Client{os.Getenv("LIBRATO_EMAIL"), os.Getenv("LIBRATO_KEY")}
  metrics  := &librato.Metrics{Source: source, Gauges: data}
  response := client.PostMetrics(metrics)

  if response != nil {
    panic(response)
  } else {
    fmt.Printf("Recorded stats to Librato: %v\n", data)
  }
}

func main() {
  stats := Collect(os.Getenv("RABBIT_QUEUE"))
  Record(stats, os.Getenv("LIBRATO_SOURCE"))
}
