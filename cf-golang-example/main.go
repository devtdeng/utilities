package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func largeResponseHandler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, string("Response from POST"))
			case http.MethodGet:
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, randStringBytes(10240000))
			}
		})
}

func delayResponseHandler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Second)
			switch r.Method {
			case http.MethodPost:
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, string("Response from POST"))
			case http.MethodGet:
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, string("Response from GET"))
			}
		})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "This is home page!")
	})
	http.HandleFunc("/large", largeResponseHandler())
	http.HandleFunc("/delay", delayResponseHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
