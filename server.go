package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

const maxMessageSize = 256

type response struct {
	Source       string `json:"source"`
	Message      string `json:"message"`
	ResponseTime int64  `json:"responseTime"`
}

func NewMessageServer(path, redisAddr string) (http.Handler, error) {
	pool := redis.NewPool(func() (redis.Conn, error) {
		return redis.Dial("tcp", redisAddr)
	}, 10)
	var (
		r = mux.NewRouter()
		m = &MessageServer{
			r:        r,
			dataPath: path,
			pool:     pool,
		}
	)
	r.HandleFunc("/{key:.*}", m.get).Methods("GET")
	r.HandleFunc("/{key:.*}", m.post).Methods("POST")
	return m, nil
}

type MessageServer struct {
	r        *mux.Router
	dataPath string
	pool     *redis.Pool
}

func (m *MessageServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("new request")
	m.r.ServeHTTP(w, r)
}

func (m *MessageServer) get(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var (
		source = "cache"
		key    = mux.Vars(r)["key"]
	)
	msg, err := m.fetchMessageFromCache(key)
	if err != nil {
		if err != errIsNotExist {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"key":   key,
			}).Error("fetch message from cache")
			http.Error(w, "fetch from cache", http.StatusInternalServerError)
			return
		}
		source = "disk"
		if msg, err = m.fetchMessageFromDisk(key); err != nil {
			if err == errIsNotExist {
				http.Error(w, "key does not exist", http.StatusNotFound)
				return
			}
			logrus.WithFields(logrus.Fields{
				"error": err,
				"key":   key,
			}).Error("fetch message from disk")
			http.Error(w, "fetch from disk", http.StatusInternalServerError)
			return
		}
		if err := m.saveMessageInCache(key, msg); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"key":   key,
			}).Error("save to cache")
			http.Error(w, "save to cache", http.StatusInternalServerError)
			return
		}
	}
	resp := response{
		Source:       source,
		Message:      msg,
		ResponseTime: time.Now().Sub(start).Nanoseconds() / 1000000,
	}
	logrus.WithFields(logrus.Fields{
		"responseTime": resp.ResponseTime,
		"key":          key,
		"source":       source,
	}).Debug("response")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"key":   key,
		}).Error("marshal response")
	}
}

func (m *MessageServer) post(w http.ResponseWriter, r *http.Request) {
	if err := m.validateRequest(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	key := mux.Vars(r)["key"]
	// make it so keys are unique and can only be created not updated.
	f, err := os.OpenFile(filepath.Join(m.dataPath, key), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"key":   key,
		}).Error("create key file")
		http.Error(w, "invalid message", http.StatusInternalServerError)
		return
	}
	defer f.Close()
	if _, err := io.Copy(f, r.Body); err != nil {
		logrus.WithField("error", err).Error("write body")
		http.Error(w, "write message", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (m *MessageServer) fetchMessageFromCache(key string) (string, error) {
	msg, err := redis.String(m.do("GET", key))
	if err != nil {
		if err == redis.ErrNil {
			return "", errIsNotExist
		}
		return "", err
	}
	return msg, nil
}

func (m *MessageServer) fetchMessageFromDisk(key string) (string, error) {
	f, err := os.Open(filepath.Join(m.dataPath, key))
	if err != nil {
		if os.IsNotExist(err) {
			return "", errIsNotExist
		}
		return "", err
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	// cheatercheatercheatercheatercheatercheatercheatercheatercheater
	// cheatercheatercheatercheatercheatercheatercheatercheatercheater
	// cheatercheatercheatercheatercheatercheatercheatercheatercheater
	// cheatercheatercheatercheatercheatercheatercheatercheatercheater
	// cheatercheatercheatercheatercheatercheatercheatercheatercheater
	// cheatercheatercheatercheatercheatercheatercheatercheatercheater \
	time.Sleep(2 * time.Second)
	// cheatercheatercheatercheatercheatercheatercheatercheatercheater /
	// cheatercheatercheatercheatercheatercheatercheatercheatercheater
	// cheatercheatercheatercheatercheatercheatercheatercheatercheater
	// cheatercheatercheatercheatercheatercheatercheatercheatercheater
	// cheatercheatercheatercheatercheatercheatercheatercheatercheater
	// cheatercheatercheatercheatercheatercheatercheatercheatercheater
	return string(data), nil
}

func (m *MessageServer) saveMessageInCache(key, msg string) error {
	_, err := m.do("SET", key, msg)
	return err
}

func (m *MessageServer) do(cmd string, args ...interface{}) (interface{}, error) {
	conn := m.pool.Get()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

// validateRequest ensures that the message body is no larger than the specified
// maximum size.
func (m *MessageServer) validateRequest(r *http.Request) error {
	if r.ContentLength > maxMessageSize {
		return fmt.Errorf("message size cannot be larger than %d", maxMessageSize)
	}
	if r.ContentLength == 0 {
		return fmt.Errorf("message size must have content")
	}
	return nil
}
