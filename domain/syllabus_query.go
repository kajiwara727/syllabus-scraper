package domain

type SyllabusQuery struct {
	Faculty string
	Year    string
	Term    *string
	Week    []string
	Period  []string
}
