package handlers

import (
	"app/services"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func SyncHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value, err := services.SyncValue()
	if err != nil {
		log.Printf("Sync failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"valueSynced": value,
	})
}
