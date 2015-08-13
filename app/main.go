package main

import (
	"errors"
	"flag"
	"net/http"

	"github.com/Sirupsen/logrus"
)

var errIsNotExist = errors.New("key does not exist")

var (
	addr         string
	redisAddress string
	debug        bool
)

func init() {
	flag.StringVar(&addr, "addr", "localhost:8123", "address for the server to bind to")
	flag.StringVar(&redisAddress, "redis", "localhost:6379", "redis server address")
	flag.BoolVar(&debug, "debug", false, "enable debug output in the logs")
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
