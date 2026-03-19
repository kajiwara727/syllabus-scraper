package infrastructure

import (
	"encoding/json"
	"os"
	"strings"
	"syllabus-scraper/domain"
)

type SyllabusJSONRepository struct {
	allData []domain.Syllabus
}

// 起動時にJSONファイルを読み込んでメモリにキャッシュする
func NewSyllabusJSONRepository(filePath string) (*SyllabusJSONRepository, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []domain.Syllabus
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	return &SyllabusJSONRepository{
		allData: data,
	}, nil
}

// メモリ上のデータから条件に合致するものを検索して返す
func (r *SyllabusJSONRepository) GetSyllabus(query domain.SyllabusQuery) ([]domain.Syllabus, error) {
	var results []domain.Syllabus

	for _, s := range r.allData {
		// ここで query の条件に基づくフィルタリングを実装します。
		// 例: CampusInfo や WeekDayPeriod などの文字列マッチング
		// (Aura API側で処理されていたフィルタリングを独自に実装する必要があります)

		match := true

		// 曜日・時限のフィルタリングの例 (実際のデータフォーマットに合わせて調整してください)
		if len(query.Week) > 0 {
			weekMatch := false
			for _, w := range query.Week {
				if strings.Contains(s.WeekDayPeriod, w) {
					weekMatch = true
					break
				}
			}
			if !weekMatch {
				match = false
			}
		}

		if match {
			results = append(results, s)
		}
	}

	return results, nil
}
