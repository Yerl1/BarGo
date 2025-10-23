package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve from /app/web instead of /app
	fs := http.FileServer(http.Dir("/app/web"))
	http.Handle("/", fs)

	log.Printf("🚀 BarGo server running on port %s...", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
