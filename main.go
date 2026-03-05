package main

import (
	"fmt"

	"syllabus-scraper/domain"
	"syllabus-scraper/infrastructure"
	"syllabus-scraper/usecase"
)

func main() {

	repo := infrastructure.NewSyllabusAPI()

	uc := usecase.NewSyllabusUsecase(repo)

	term := "春学期"

	query := domain.SyllabusQuery{
		Faculty: "21",          //
		Year:    "2026",        // "2024~"
		Term:    &term,         // "春学期", "秋学期"
		Week:    []string{"月"}, // "月~金"
		Period:  []string{""},  // "1~6"
	}

	result, err := uc.GetSyllabus(query)

	if err != nil {
		panic(err)
	}

	for _, s := range result {
		fmt.Println(s.ID, s.CourseName, s.PersonalName)
	}
}
