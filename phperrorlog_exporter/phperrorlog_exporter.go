package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/foomo/phperrorlog_exporter/logparser"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	flagHTTP := flag.String("http", ":9423", "Address to listen on for web interface.")
	flag.Parse()
	usage := func() {
		fmt.Println("Usage:", os.Args[0], "domain:path/to/php/error.log [domain:path/to/other/log]")
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

	phpErrors := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "php_errors",
			Help: "How many PHP errors partitioned by type and site",
		},
		[]string{"type", "site"},
	)
	prometheus.MustRegister(phpErrors)

	if len(*flagHTTP) > 0 {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			http.ListenAndServe(*flagHTTP, nil)
		}()
        fmt.Println("Listening on address:port => ", *flagHTTP)
	}

	for {
		select {
		case observation := <-chanObservation:
			fmt.Println("observing", observation.Name, "stats", observation.Stats)
			for t, v := range observation.Stats {
				phpErrors.WithLabelValues(t, observation.Name).Observe(float64(v))
			}
		}
	}
}
