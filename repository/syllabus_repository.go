package repository

import (
	"syllabus-scraper/domain"
)

type SyllabusRepository interface {
	GetSyllabus(query domain.SyllabusQuery) ([]domain.Syllabus, error)
}
