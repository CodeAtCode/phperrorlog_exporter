package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/phperrorlog_exporter/logparser"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	flagHTTP := flag.String("http", "", "expose metrics ip:port")
	flag.Parse()
	usage := func() {
		fmt.Println("Usage:", os.Args[0], "name:path/to/php/error.log [name:path/to/other/log]")
		flag.PrintDefaults()
		os.Exit(1)

	}

	args := flag.Args()

	if len(args) < 1 {
		usage()
	}
	chanObservation := make(chan logparser.Observation)
	for _, nameFile := range args {
		parts := strings.Split(nameFile, ":")
		if len(parts) != 2 {
			usage()
		}
		go logparser.Observe(parts[0], parts[1], chanObservation, time.Second*10)
	}

	phpErrors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "php_errors",
			Help: "How many PHP errors partitioned by type and site",
		},
		[]string{"type", "site"},
	)
	prometheus.MustRegister(phpErrors)

	if len(*flagHTTP) > 0 {
		go func() {
			http.Handle("/metrics", prometheus.Handler())
			http.ListenAndServe(*flagHTTP, nil)
		}()
	}

	for {
		select {
		case observation := <-chanObservation:
			fmt.Println("observing", observation.Name, "stats", observation.Stats)
			for t, v := range observation.Stats {
				phpErrors.WithLabelValues(t, observation.Name).Set(float64(v))
			}
		}
	}
}
