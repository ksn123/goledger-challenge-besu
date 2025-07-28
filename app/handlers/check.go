package handlers

import (
	"app/services"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func CheckHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	match, err := services.CompareContractWithDB()
	if err != nil {
		log.Printf("Comparison failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"equal": match,
	})
}
