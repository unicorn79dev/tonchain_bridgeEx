package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// DB Connection Pool
var db *sql.DB

// Initialize MySQL connection
func initDB() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/tonbridge" // Update with your DB details
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}
}

func getClientIP(r *http.Request) string {
	// Check for X-Forwarded-For header
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For can have multiple IPs; the first one is the client IP
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}

	// Check for X-Real-IP header
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Fallback to RemoteAddr
	return r.RemoteAddr
}

// Save request data to MySQL
func saveRequestToDB(r *http.Request) {
	IP := getClientIP(r)
	Origin := r.Header.Get("Origin")
	Refer := r.Header.Get("Referer")
	UserAgent := r.UserAgent()

	// Insert the data into MySQL
	query := "INSERT INTO connectlog (origin, refer, ip, user_agent) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, Origin, Refer, IP, UserAgent)
	if err != nil {
		log.Printf("Error saving request to DB: %v", err)
	} else {
		fmt.Println("Request logged to DB")
	}
}
