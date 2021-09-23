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
	}

	var wait sync.WaitGroup
	wait.Add(testSize)
	Limiter := rate.NewLimiter(rate.Limit(testSize), int(testSize))
	//Limiter := ratelimit.NewBucketWithQuantum(time.Second, int64(testSize), int64(testSize))
	//Limiter := ratelimit.New(int(testSize), time.Second)


	start := time.Now()
	for i:=1; i<= testSize; i++{
		go func(){
			defer wait.Done()
			if limiterFlag {
				resultBool := Limiter.Allow()
				//_, resultBool := Limiter.TakeMaxDuration(1, 0)
				//resultBool := !Limiter.Limit()
				if !resultBool{
					fmt.Println("Fuck you")
					return
				}

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
	//fmt.Println(duration, "fuck you")
}
