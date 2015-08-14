package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/rcrowley/go-metrics"
)

var errIsNotExist = errors.New("key does not exist")

var (
	addr         string
	redisAddress string
	debug        bool
)

func init() {
	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "address for the server to bind to")
	flag.StringVar(&redisAddress, "redis", "localhost:6379", "redis server address")
	flag.BoolVar(&debug, "debug", false, "enable debug output in the logs")
	requests = metrics.NewTimer()
	metrics.Register("requests", requests)
	go metrics.Log(metrics.DefaultRegistry, 60e9, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
}

func main() {
	flag.Parse()
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	m, err := NewMessageServer(redisAddress)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := http.ListenAndServe(addr, m); err != nil {
		logrus.Fatal(err)
	}
}
