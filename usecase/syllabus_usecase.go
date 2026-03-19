package usecase

import (
	"syllabus-scraper/domain"
	"syllabus-scraper/repository"
)

// シラバスに関するビジネスロジック
type SyllabusUsecase struct {
	repo repository.SyllabusRepository
}

func NewSyllabusUsecase(r repository.SyllabusRepository) *SyllabusUsecase {
	return &SyllabusUsecase{
		repo: r,
	}
}

// 検索条件を受け取りリポジトリ経由でシラバスデータを取得して返す
func (u *SyllabusUsecase) GetSyllabus(query domain.SyllabusQuery) ([]domain.Syllabus, error) {
	return u.repo.GetSyllabus(query)
}
