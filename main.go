package main

import (
	"encoding/json"
	"log"
	"net/http"

	"syllabus-scraper/domain"
	"syllabus-scraper/infrastructure"
	"syllabus-scraper/usecase"
)

func main() {
	// 外部APIの代わりにJSONファイルを読み込むリポジトリを使用
	repo, err := infrastructure.NewSyllabusJSONRepository("all_syllabuses.json")
	if err != nil {
		log.Fatalf("Failed to load JSON data: %v", err)
	}
	uc := usecase.NewSyllabusUsecase(repo)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Connection successful"))
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		var query domain.SyllabusQuery
		if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := uc.GetSyllabus(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
