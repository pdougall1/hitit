package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	var url string
	var channels []<-chan float64
	var responseTimes []float64
	var requestRate int
	var requestsMade int
	var timeSeconds int

	flag.IntVar(&requestRate, "r", 10, "Specify requests per second.  Defaults to 10.") // rate
	flag.StringVar(&url, "u", "http://example.com/", "Specify url of requests to make.  Defaults to http://example.com/.")
	flag.IntVar(&timeSeconds, "s", 1, "Specify number of seconds to spread requests over.  Defaults to 1.")
	flag.Parse()
	timeMilliseconds := timeSeconds * 1000

	for i := 0; i < requestRate*timeSeconds; i++ {
		channels = append(channels, makeRequest(url, timeMilliseconds))
	}

	fmt.Printf(`
__________________________________

  Request URL:  %s
  Request rate: %d
  Time:         %d

`, url, requestRate, timeSeconds)

	for responseTime := range merge(channels) {
		responseTimes = append(responseTimes, responseTime)
		requestsMade++
		fmt.Printf("\r  Successful requests #: %d", requestsMade)
	}

	printResults(GetStats(responseTimes))
}

func printResults(stats statistics) {
	response := `
__________________________________

  Mean:               %f
  Standard Deviation: %f
  Quickest:           %f
  Slowest:            %f
__________________________________

`
	mean := toSeconds(stats.mean)
	standardDeviation := toSeconds(stats.stdDev)
	quickest := toSeconds(stats.quickest)
	slowest := toSeconds(stats.slowest)
	fmt.Printf(response, mean, standardDeviation, quickest, slowest)
}

func toSeconds(nanoseconds float64) float64 {
	return nanoseconds / 1000000000
}

func makeRequest(url string, timeMilliseconds int) <-chan float64 {
	out := make(chan float64)
	go func() {

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		time.Sleep(time.Duration(r.Intn(timeMilliseconds)) * time.Millisecond)

		startTime := time.Now()
		_, err := http.Get(url)
		if err == nil {
			out <- float64(time.Since(startTime))
		}
		close(out)
	}()
	return out
}

func merge(cs []<-chan float64) <-chan float64 {
	var wg sync.WaitGroup
	out := make(chan float64)

	output := func(c <-chan float64) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
