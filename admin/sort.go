package main

import (
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

type (
	server struct {
		Name         string  `json:"name"`
		Fill         float64 `json:"fill"`
		ResponseTime float64 `json:"responseTime"`
	}
	state struct {
		TotalRequests float64  `json:"totalRequests"`
		Servers       []server `json:"servers"`
	}
)

func (s server) number() int {
	n := strings.Replace(s.Name, "frontend-", "", -1)
	i, err := strconv.Atoi(n)
	if err != nil {
		logrus.Error(err)
	}
	return i
}

type serverSorter []server

func (a serverSorter) Len() int           { return len(a) }
func (a serverSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a serverSorter) Less(i, j int) bool { return a[i].number() < a[j].number() }
