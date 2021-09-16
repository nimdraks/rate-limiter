package main

import (
	"flag"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

func main(){

	testSizeP := flag.Int("t", 200, "number of request")
	LimiterFlagP := flag.Bool("l", false, "Limiter Flag")
	flag.Parse()
	testSize := *testSizeP
	limiterFlag := *LimiterFlagP

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if _, err := io.WriteString(w, "I love you"); err != nil{
		}
	}))
	defer ts.Close()

	ServerURL := ts.URL
	httpClient := http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1000,
			MaxIdleConns: 1000,
			MaxConnsPerHost: 1000,
			IdleConnTimeout: 5 * time.Second,
		},
	}

	var wait sync.WaitGroup
	wait.Add(testSize)
	Limiter := rate.NewLimiter(rate.Limit(testSize), int(testSize))
	//fmt.Println(testSize)

	start := time.Now()
	for i:=1; i<= testSize; i++{
		go func(){
			defer wait.Done()
			if limiterFlag {
				Limiter.Allow()
			}
			body, err := httpClient.Get(ServerURL)
			if err != nil{
				fmt.Println(err)
				return
			}
			body.Body.Close()
		}()
	}

	wait.Wait()

	duration := time.Since(start)
	fmt.Println(float64(duration/time.Microsecond)/1000.0)
}
