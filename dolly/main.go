package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type server struct {
	Name          string  `json:"name"`
	Fill          float64 `json:"fill"`
	StatsEndpoint string  `json:"statsEndpoint"`
	ip            string
}

type state struct {
	TotalRequests float64  `json:"totalRequests"`
	Servers       []server `json:"servers"`
}

// getState returns the entire state of the cluster.
func getState(w http.ResponseWriter, r *http.Request) {

}

// startServer starts the web server instance with a cold cache.
func startServer(w http.ResponseWriter, r *http.Request) {

}

// cloneServer clones the server one and starts it somewhere else.
func cloneServer(w http.ResponseWriter, r *http.Request) {

}

func main() {
	s := state{
		TotalRequests: 2020,
		Servers: []server{
			{
				Name:          "linuxcon1",
				Fill:          66.6,
				StatsEndpoint: "http://127.0.0.3001",
			},
			{
				Name:          "linuxcon1",
				Fill:          100.0,
				StatsEndpoint: "http://127.0.0.3002",
			},
		},
	}
	json.NewEncoder(os.Stdout).Encode(s)
}
