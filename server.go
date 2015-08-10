package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

type response struct {
	// Fill is the % of the cache that is currently full
	Fill float64 `json:"fill"`
	// Message is the message requested by the user
	Message string `json:"message"`
}

func NewMessageServer(size int, redisAddr string) (http.Handler, error) {
	m := &MessageServer{
		r:         mux.NewRouter(),
		cacheSize: float64(size),
	}
	m.r.HandleFunc("/", m.getMessage).Methods("GET")
	m.r.HandleFunc("/cache", m.getCache).Methods("GET")
	// start filling the cache async
	go m.fillCache(redisAddr)
	return m, nil
}

type MessageServer struct {
	r         *mux.Router
	cache     []string
	cacheSize float64
	cacheLock sync.Mutex
}

func (m *MessageServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.r.ServeHTTP(w, r)
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
	m.cacheLock.Lock()
	n := len(m.cache)
	if n == 0 {
		m.cacheLock.Unlock()
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}
	var resp response
	resp.Message = m.cache[rand.Intn(n)]
	m.cacheLock.Unlock()
	resp.Fill = float64(float64(n)/m.cacheSize) * 100.0
	// sleep to simulate slow response times while cache is still filling
	sleepTime := time.Duration(10.0 - (10.0 * float64(float64(n)/m.cacheSize)))
	time.Sleep(sleepTime * time.Second)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logrus.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (m *MessageServer) fillCache(redisAddr string) {
	for i := 0; i < 5; i++ {
		conn, err := redis.Dial("tcp", redisAddr)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		defer conn.Close()
		data, err := redis.Strings(conn.Do("LRANGE", "messages", 0, -1))
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		for _, d := range data {
			m.cacheLock.Lock()
			m.cache = append(m.cache, d)
			m.cacheLock.Unlock()
			time.Sleep(2 * time.Second)
		}
		return
	}
	if len(m.cache) == 0 {
		logrus.Fatalf("unable to fill cache from redis @ %q", redisAddr)
	}
}
