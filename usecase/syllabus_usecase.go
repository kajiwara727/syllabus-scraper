package usecase

import (
	"syllabus-scraper/domain"
	"syllabus-scraper/repository"
)

type SyllabusUsecase struct {
	repo repository.SyllabusRepository
}

func NewSyllabusUsecase(r repository.SyllabusRepository) *SyllabusUsecase {
	return &SyllabusUsecase{
		repo: r,
	}
}

func (u *SyllabusUsecase) GetSyllabus(query domain.SyllabusQuery) ([]domain.Syllabus, error) {
	return u.repo.GetSyllabus(query)
}
