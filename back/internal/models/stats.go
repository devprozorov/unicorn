package models

// HomeStats содержит статистику для главной страницы
type HomeStats struct {
	TotalResumes   int64 `json:"total_resumes"`
	TotalUsers     int64 `json:"total_users"`
	TotalCompanies int64 `json:"total_companies"`
	TotalVacancies int64 `json:"total_vacancies"`
	CachedAt       int64 `json:"cached_at"` // Unix timestamp
}
