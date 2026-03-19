package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"syllabus-scraper/domain"
	"syllabus-scraper/infrastructure"
)

func main() {
	api := infrastructure.NewSyllabusAPI() // 既存のAPIクライアントを利用

	var allSyllabuses []domain.Syllabus

	// 実際には、"limits"を大きくするか、
	// 学部や曜日ごとにループを回してすべてのデータを取得する必要があります。
	// ここでは例として、特定の検索条件で取得したものを蓄積するイメージです。
	faculties := []string{"59"} // 例

	for _, faculty := range faculties {
		query := domain.SyllabusQuery{
			Faculty: faculty,
			Year:    "2026", // 今年度など
		}

		results, err := api.GetSyllabus(query)
		if err != nil {
			log.Printf("Error fetching %s: %v", faculty, err)
			continue
		}
		allSyllabuses = append(allSyllabuses, results...)
	}

	// JSONファイルとして保存
	file, err := os.Create("59.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(allSyllabuses); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total %d records saved to 11.json\n", len(allSyllabuses))
}
