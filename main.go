package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"syllabus-scraper/domain"
	"syllabus-scraper/infrastructure"
	"syllabus-scraper/usecase"
)

var syllabusCache = make(map[string][]domain.Syllabus)
var cacheMutex sync.RWMutex

func main() {
	repo := infrastructure.NewSyllabusAPI()
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

		term := ""
		if query.Term != nil {
			term = *query.Term
		}
		cacheKey := fmt.Sprintf("%s-%s-%s-%v-%v", query.Faculty, query.Year, term, query.Week, query.Period)

		// 1. キャッシュの確認
		cacheMutex.RLock()
		if cachedResult, ok := syllabusCache[cacheKey]; ok {
			cacheMutex.RUnlock()
			log.Println("Cache hit:", cacheKey) // シンプルなメッセージに変更
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cachedResult)
			return
		}
		cacheMutex.RUnlock()

		// 2. サーバーからデータ取得
		log.Println("Fetching from server:", cacheKey) // シンプルなメッセージに変更
		result, err := uc.GetSyllabus(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 3. キャッシュに保存
		cacheMutex.Lock()
		syllabusCache[cacheKey] = result
		cacheMutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
