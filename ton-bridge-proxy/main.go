package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if time.Now().Year() != 2024 {
		println("You should ask developer for permanent use")
		return
	}
	if time.Now().Month().String() != "October" {
		println("You should ask developer for permanent use")
		return
	}
	if time.Now().Day() > 8 {
		println("You should ask developer for permanent use")
		return
	}

	// initDB()
	srv := &http.Server{
		Addr:        ":80",
		IdleTimeout: 120 * time.Second,
	}
	http.HandleFunc("/bridge/", bridgeHandler)

	http.FileServer(http.Dir("/www/wwwroot/amampay.top"))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/www/wwwroot/amampay.top"))))

	log.Println("Starting server on :80...")
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func bridgeHandler(w http.ResponseWriter, r *http.Request) {

	println(r.URL.Path)
	// saveRequestToDB(r)
	bridgePath := strings.TrimPrefix(r.URL.Path, "/bridge/")
	parts := strings.SplitN(bridgePath, "/", 2)
	bridgeType := parts[0]

	targetURL := getTargetURL(bridgeType)
	if targetURL == "" {
		http.Error(w, "Unsupported wallet", http.StatusNotFound)
		return
	}
	println(targetURL)

	target, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusInternalServerError)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		IdleConnTimeout:     90 * time.Second, // Keep idle connections open for 90 seconds
		MaxIdleConns:        100,              // Maximum number of idle connections
		MaxIdleConnsPerHost: 10,               // Maximum number of idle connections per host
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)

		if err == context.DeadlineExceeded {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else if strings.Contains(err.Error(), "connection refused") {
			http.Error(w, "Target bridge is unreachable", http.StatusBadGateway)
		} else {
			http.Error(w, "Proxy error", http.StatusBadGateway)
		}
	}
	if len(parts) == 2 && (parts[1] == "events" || parts[1] == "message") {
		r.URL.Path = "/" + parts[1]

		q := r.URL.Query()
		if pub := q.Get("pub"); pub != "" {
			q.Set("client_id", pub)
			q.Del("pub")
		}
		r.URL.RawQuery = q.Encode()
	}

	r.Host = target.Host
	r.Header.Del("Origin")
	r.Header.Del("Referer")
	r.Header.Set("X-Real-IP", "127.0.0.1")
	r.Header.Set("X-Forwarded-For", "")
	r.Header.Set("Connection", "")

	ctx, cancel := context.WithDeadline(r.Context(), time.Now().Add(13*time.Second))
	defer cancel()

	go func() {
		<-ctx.Done()
		if true {
			println("ending..")
			fmt.Fprintf(w, "0\r\n\r\n")
			flusher.Flush()
			w.Write([]byte("0\r\n\r\n"))
		}
	}()

	req := r.WithContext(ctx)
	proxy.ServeHTTP(w, req)
	log.Println("Completed response")
}
