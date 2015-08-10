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
	cacheSize    int
	debug        bool
)

func init() {
	flag.StringVar(&addr, "addr", "localhost:8123", "address for the server to bind to")
	flag.StringVar(&redisAddress, "redis", "localhost:6379", "redis server address")
	flag.BoolVar(&debug, "debug", false, "enable debug output in the logs")
	flag.IntVar(&cacheSize, "size", -1, "number of elements in the cache")
}

func main() {
	flag.Parse()
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if cacheSize == -1 {
		logrus.Fatal("-cacheSize must be specified")
	}
	m, err := NewMessageServer(cacheSize, redisAddress)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := http.ListenAndServe(addr, m); err != nil {
		logrus.Fatal(err)
	}
}
