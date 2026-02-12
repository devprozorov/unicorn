package handlers

import (
	"net/http"
	"time"

	"unicorn-auth/internal/models"
	"unicorn-auth/internal/repo"
	"unicorn-auth/internal/security"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	users     *repo.UserRepo
	resumes   *repo.ResumeRepo
	vacancies *repo.VacancyRepo
	cacher    *security.StatsCacher
}

func NewHomeHandler(users *repo.UserRepo, resumes *repo.ResumeRepo, vacancies *repo.VacancyRepo) *HomeHandler {
	return &HomeHandler{
		users:     users,
		resumes:   resumes,
		vacancies: vacancies,
		cacher:    security.NewStatsCacherWithDuration(30 * time.Second),
	}
}

// GetStats возвращает статистику главной страницы с кэшированием
func (h *HomeHandler) GetStats(c *gin.Context) {
	// Проверяем кэш
	if cached := h.cacher.Get(); cached != nil {
		c.JSON(http.StatusOK, cached)
		return
	}

	ctx := c.Request.Context()

	// Получаем все счетчики параллельно (для оптимизации)
	var totalResumes, totalVacancies, totalUsers, totalCompanies int64
	var resumeErr, vacancyErr, userErr, companyErr error

	// Используем горутины для параллельного получения данных
	done := make(chan struct{})

	go func() {
		totalResumes, resumeErr = h.resumes.CountActive(ctx)
		done <- struct{}{}
	}()

	go func() {
		totalVacancies, vacancyErr = h.vacancies.CountActive(ctx)
		done <- struct{}{}
	}()

	go func() {
		totalUsers, userErr = h.users.CountAll(ctx)
		done <- struct{}{}
	}()

	go func() {
		totalCompanies, companyErr = h.users.CountByType(ctx, "company")
		done <- struct{}{}
	}()

	// Ждем завершения всех горутин
	for i := 0; i < 4; i++ {
		<-done
	}

	// Проверяем ошибки
	if resumeErr != nil || vacancyErr != nil || userErr != nil || companyErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch statistics"})
		return
	}

	stats := &models.HomeStats{
		TotalResumes:   totalResumes,
		TotalVacancies: totalVacancies,
		TotalUsers:     totalUsers,
		TotalCompanies: totalCompanies,
		CachedAt:       time.Now().UTC().Unix(),
	}

	// Кэшируем результат
	h.cacher.Set(stats)

	c.JSON(http.StatusOK, stats)
}
