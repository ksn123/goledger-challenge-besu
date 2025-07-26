package main

import (
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "app/handlers"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: No .env file found")
    }

    http.HandleFunc("/get", handlers.GetHandler)
    http.HandleFunc("/set", handlers.SetHandler)
    http.HandleFunc("/sync", handlers.SyncHandler)
    http.HandleFunc("/check", handlers.CheckHandler)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Server running on port %s...", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
