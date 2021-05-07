package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/CodeAtCode/phperrorlog_exporter/logparser"
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
	files := 0
	chanObservation := make(chan logparser.Observation)
	for _, nameFile := range args {
		parts := strings.Split(nameFile, ":")
		if len(parts) != 2 {
			usage()
		}
		files++
		
		go logparser.Observe(parts[0], parts[1], chanObservation, time.Second*10)
	}

	phpErrors := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "php_errors",
			Help: "How many PHP errors partitioned by type and site",
            MaxAge: 1,
            AgeBuckets: 1,
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

	run := 0
	for {
		select {
		case observation := <-chanObservation:
            if run == files {
                run = 0
                phpErrors.Reset()
                fmt.Println("Reset done")
            }
			fmt.Println("observing", observation.Name, "stats", observation.Stats)
			for t, v := range observation.Stats {
				phpErrors.WithLabelValues(t, observation.Name).Observe(float64(v))
			}
            run++
		}
	}
}
