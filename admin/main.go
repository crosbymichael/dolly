package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"sort"

	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
)

var redisAddr string

func init() {
	flag.StringVar(&redisAddr, "redis", "", "redis address for config")
	flag.Parse()
}

var pool *redis.Pool = redis.NewPool(func() (redis.Conn, error) {
	return redis.Dial("tcp", redisAddr)
}, 10)

func do(cmd string, args ...interface{}) (interface{}, error) {
	conn := pool.Get()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

// getState returns the entire state of the cluster.
func getState(w http.ResponseWriter, r *http.Request) {
	writeCORS(w)
	req, err := getTotalRequests()
	if err != nil {
		logrus.Error(err)
	}
	s := state{
		TotalRequests: float64(req),
	}
	for _, n := range nodes {
		s.Servers = append(s.Servers, server{
			Name:         n.name,
			Fill:         n.getFill(),
			ResponseTime: n.getResponseTime(),
		})
	}
	ss := serverSorter(s.Servers)
	sort.Sort(ss)
	// SORT servers array
	if err := json.NewEncoder(w).Encode(s); err != nil {
		logrus.Error(err)
	}
}

func getTotalRequests() (int64, error) {
	conn := pool.Get()
	defer conn.Close()
	return redis.Int64(conn.Do("GET", "requests"))
}

// startServer starts the web server instance with a cold cache.
func startServer(w http.ResponseWriter, r *http.Request) {
	writeCORS(w)
	name := mux.Vars(r)["name"]
	n := nodes[name]
	if err := n.start(false); err != nil {
		logrus.Error(err)
	}
}

func stopServer(w http.ResponseWriter, r *http.Request) {
	writeCORS(w)
	name := mux.Vars(r)["name"]
	n := nodes[name]
	if err := n.stop(); err != nil {
		logrus.Error(err)
	}
}

// cloneServer clones the server one and starts it somewhere else.
func cloneServer(w http.ResponseWriter, r *http.Request) {
	writeCORS(w)
	r.ParseForm()
	//	n1 := nodes[mux.Vars(r)["name"]]
	//	n2 := ndoes[r.Form.Get("server")]
}

func httpError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func writeCORS(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-Registry-Auth")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
}

func main() {
	if err := loadNodes(); err != nil {
		logrus.Fatal(err)
	}
	// router setup
	r := mux.NewRouter()
	r.HandleFunc("/", getState).Methods("GET")
	r.HandleFunc("/{name:.*}/start", startServer).Methods("POST")
	r.HandleFunc("/{name:.*}/stop", stopServer).Methods("POST")
	r.HandleFunc("/{name:.*}/clone", cloneServer).Methods("POST")
	r.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		writeCORS(w)
	}).Methods("OPTIONS")

	// start the server
	if err := http.ListenAndServe("127.0.0.1:8765", r); err != nil {
		logrus.Fatal(err)
	}
}
