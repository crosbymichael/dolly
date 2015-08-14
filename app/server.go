package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/rcrowley/go-metrics"
)

var requests metrics.Timer

type response struct {
	// Fill is the % of the cache that is currently full
	Fill float64 `json:"fill"`
	// Message is the message requested by the user
	Message string `json:"message"`
}

func NewMessageServer(redisAddr string) (http.Handler, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	m := &MessageServer{
		r:    mux.NewRouter(),
		name: strings.Replace(hostname, "linuxcon-demo-", "", -1),
		pool: redis.NewPool(func() (redis.Conn, error) {
			return redis.Dial("tcp", redisAddr)
		}, 10),
	}
	m.r.HandleFunc("/", m.getMessage).Methods("GET")
	m.r.HandleFunc("/cache", m.getCache).Methods("GET")
	// start filling the cache async on boot
	go m.fillCache(redisAddr)
	go func() {
		for range time.Tick(5 * time.Second) {
			if _, err := m.do("SET", fmt.Sprintf("nodes.%s.avg", m.name), requests.RateMean()); err != nil {
				logrus.Error(err)
			}
		}
	}()
	return m, nil
}

type MessageServer struct {
	r         *mux.Router
	cache     []string
	cacheSize float64
	cacheLock sync.Mutex
	pool      *redis.Pool
	name      string
}

func (m *MessageServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn := m.pool.Get()
	conn.Do("INCR", "requests")
	conn.Close()
	start := time.Now()
	m.r.ServeHTTP(w, r)
	requests.Update(time.Now().Sub(start))
}

func (m *MessageServer) getCache(w http.ResponseWriter, r *http.Request) {
	m.cacheLock.Lock()
	n := len(m.cache)
	m.cacheLock.Unlock()
	if err := json.NewEncoder(w).Encode(struct {
		Fill float64 `json:"fill"`
	}{
		Fill: (float64(n) / m.cacheSize) * 100.0,
	}); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (m *MessageServer) getMessage(w http.ResponseWriter, r *http.Request) {
	// sleep to simulate slow response times while cache is still filling
	fill, msg, err := m.fetchNewMessage()
	if err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}
	resp := response{
		Fill:    fill,
		Message: msg,
	}
	// sleep to simulate slow response times while cache is filling
	sleepTime := time.Duration(10.0 - (10.0 * (fill / 100.0)))
	if sleepTime.Seconds() != 0 {
		time.Sleep(sleepTime * time.Second)
	}
	// send the request
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (m *MessageServer) fillCache(redisAddr string) {
	f, err := os.Open("/data.json")
	if err != nil {
		logrus.Fatal(err)
	}
	var data []string
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		logrus.Fatal(err)
	}
	f.Close()
	m.cacheSize = float64(len(data))
	for i, d := range data {
		m.cacheLock.Lock()
		m.cache = append(m.cache, d)
		m.cacheLock.Unlock()
		time.Sleep(2 * time.Second)
		f := (float64(i+1) / m.cacheSize) * 100.0
		if _, err := m.do("SET", fmt.Sprintf("nodes.%s.fill", m.name), f); err != nil {
			logrus.Error(err)
		}
	}
}

func (m *MessageServer) do(cmd string, args ...interface{}) (interface{}, error) {
	conn := m.pool.Get()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

// fetchNewMessage returns the cache fill % along with a random message
func (m *MessageServer) fetchNewMessage() (float64, string, error) {
	m.cacheLock.Lock()
	defer m.cacheLock.Unlock()
	n := len(m.cache)
	if n == 0 {
		return 0.0, "", fmt.Errorf("no content in cache")
	}
	return float64(float64(n)/m.cacheSize) * 100.0, m.cache[rand.Intn(n)], nil
}
