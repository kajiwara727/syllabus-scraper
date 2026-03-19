package repository

import (
	"syllabus-scraper/domain"
)

// シラバスデータ取得のインタフェース
type SyllabusRepository interface {
	GetSyllabus(query domain.SyllabusQuery) ([]domain.Syllabus, error)
}
