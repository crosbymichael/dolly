package main

import (
	"errors"
	"flag"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
)

var errIsNotExist = errors.New("key does not exist")

var (
	addr         string
	dataPath     string
	redisAddress string
	debug        bool
)

func init() {
	flag.StringVar(&addr, "addr", "localhost:8123", "address for the server to bind to")
	flag.StringVar(&dataPath, "path", "/tmp/yolo", "data path to save keys to")
	flag.StringVar(&redisAddress, "redis", "localhost:6379", "redis server address")
	flag.BoolVar(&debug, "debug", false, "enable debug output in the logs")
}

func main() {
	flag.Parse()
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		logrus.Fatal(err)
	}
	m, err := NewMessageServer(dataPath, redisAddress)
	if err != nil {
		logrus.Fatal(err)
	}
	if err := http.ListenAndServe(addr, m); err != nil {
		logrus.Fatal(err)
	}
}
