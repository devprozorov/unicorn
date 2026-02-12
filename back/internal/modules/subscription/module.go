package subscription

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"unicorn-auth/internal/http/middleware"
	"unicorn-auth/internal/models"
	"unicorn-auth/internal/repo"
	"unicorn-auth/internal/security"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Config struct {
	Price            string
	DurationDays     int
	RobokassaEnabled bool
}

func Register(r *gin.Engine, cfg Config, sec *security.Security, robokassa *security.Robokassa,
	users *repo.UserRepo, subs *repo.SubscriptionRepo, vacancies *repo.VacancyRepo, resumes *repo.ResumeRepo) {

	api := r.Group("/api")

	// Защищенные эндпоинты (требуют авторизации)
	protected := api.Group("")
	protected.Use(middleware.RequireAuth(sec))

	// GET /api/subscription/status - Получить статус подписки
	protected.GET("/subscription/status", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)

		u, err := users.FindByUserID(c.Request.Context(), uid)
		if err != nil || u == nil {
			c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
			return
		}

		activeSub, _ := subs.GetActiveByUserID(c.Request.Context(), uid)

		response := gin.H{
			"ok":     true,
			"active": u.Subscription.Active,
		}

		if activeSub != nil {
			response["endDate"] = activeSub.EndDate
			response["daysLeft"] = int(time.Until(activeSub.EndDate).Hours() / 24)
		}

		c.JSON(200, response)
	})

	// POST /api/subscription/create-payment - Создать ссылку на оплату
	protected.POST("/subscription/create-payment", func(c *gin.Context) {
		uid := c.GetString(middleware.CtxUserID)

		u, err := users.FindByUserID(c.Request.Context(), uid)
		if err != nil || u == nil {
			c.JSON(401, gin.H{"ok": false, "error": "unauthorized"})
			return
		}

		if !cfg.RobokassaEnabled {
			c.JSON(503, gin.H{"ok": false, "error": "payment_disabled"})
			return
		}

		// Создаем запись подписки в статусе pending
		invID := time.Now().Unix() // Используем timestamp как InvID
		sub := &models.Subscription{
			UserID:   uid,
			Amount:   parseFloat(cfg.Price),
			Currency: "RUB",
			Status:   "pending",
			InvID:    invID,
			OutSum:   cfg.Price,
		}

		if err := subs.Create(c.Request.Context(), sub); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": "server_error"})
			return
		}

		// Генерируем ссылку на оплату Robokassa
		description := fmt.Sprintf("Подписка на %d дней", cfg.DurationDays)
		paymentURL := robokassa.GeneratePaymentURL(cfg.Price, invID, description, uid)

		log.Printf("subscription: payment created for user=%s, invID=%d, amount=%s, url=%s", uid, invID, cfg.Price, paymentURL)

		c.JSON(200, gin.H{
			"ok":         true,
			"paymentUrl": paymentURL,
			"invId":      invID,
			"amount":     cfg.Price,
		})
	})

	// POST /api/subscription/robokassa/result - Callback от Robokassa (Result URL)
	api.POST("/subscription/robokassa/result", func(c *gin.Context) {
		outSum := c.PostForm("OutSum")
		invIDStr := c.PostForm("InvId")
		signatureValue := c.PostForm("SignatureValue")
		userID := c.PostForm("Shp_userId")

		log.Printf("robokassa result: received callback - OutSum=%s, InvId=%s, Signature=%s, UserId=%s",
			outSum, invIDStr, signatureValue, userID)

		invID, err := strconv.ParseInt(invIDStr, 10, 64)
		if err != nil {
			log.Printf("robokassa result: invalid invID format: %v", err)
			c.String(400, "bad request")
			return
		}

		// Проверяем подпись
		if !robokassa.VerifyResultSignature(outSum, invID, signatureValue, userID) {
			log.Printf("robokassa result: invalid signature for invID=%d, expected signature calculation based on OutSum=%s", invID, outSum)
			c.String(400, "invalid signature")
			return
		}

		log.Printf("robokassa result: signature verified successfully for invID=%d", invID)

		// Находим подписку
		sub, err := subs.FindByInvID(c.Request.Context(), invID)
		if err != nil || sub == nil {
			log.Printf("robokassa result: subscription not found invID=%d, error=%v", invID, err)
			c.String(404, "subscription not found")
			return
		}

		log.Printf("robokassa result: found subscription %s for user %s", sub.SubscriptionID, userID)

		// Обновляем статус подписки
		startDate := time.Now().UTC()
		endDate := startDate.AddDate(0, 0, cfg.DurationDays)

		if err := subs.UpdateStatus(c.Request.Context(), invID, "paid", startDate, endDate); err != nil {
			log.Printf("robokassa result: failed to update subscription: %v", err)
			c.String(500, "server error")
			return
		}

		log.Printf("robokassa result: subscription status updated to paid, endDate=%s", endDate)

		// Обновляем пользователя
		if err := users.UpdateByUserID(c.Request.Context(), userID, bson.M{
			"subscription.active": true,
			"subscription.until":  endDate,
		}); err != nil {
			log.Printf("robokassa result: failed to update user: %v", err)
			c.String(500, "server error")
			return
		}

		log.Printf("robokassa result: user subscription activated for user=%s until=%s", userID, endDate)

		// Обновляем все вакансии и резюме пользователя - добавляем премиум статус
		colorCode := "#FFD700" // Gold color for premium
		if err := updateUserContentToPremium(c.Request.Context(), userID, colorCode, vacancies, resumes); err != nil {
			log.Printf("robokassa result: failed to update content: %v", err)
		} else {
			log.Printf("robokassa result: content updated to premium for user=%s", userID)
		}

		log.Printf("robokassa result: payment processing completed successfully for user=%s, invID=%d", userID, invID)
		c.String(200, fmt.Sprintf("OK%d", invID))
	})

	// GET /api/subscription/robokassa/success - Success URL (переадресация после оплаты)
	api.GET("/subscription/robokassa/success", func(c *gin.Context) {
		outSum := c.Query("OutSum")
		invIDStr := c.Query("InvId")
		signatureValue := c.Query("SignatureValue")
		userID := c.Query("Shp_userId")

		invID, err := strconv.ParseInt(invIDStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"ok": false, "error": "bad_request"})
			return
		}

		// Проверяем подпись
		if !robokassa.VerifySuccessSignature(outSum, invID, signatureValue, userID) {
			c.JSON(400, gin.H{"ok": false, "error": "invalid_signature"})
			return
		}

		c.JSON(200, gin.H{
			"ok":      true,
			"message": "Оплата успешно завершена!",
			"invId":   invID,
		})
	})

	// GET /api/subscription/robokassa/fail - Fail URL (переадресация при ошибке)
	api.GET("/subscription/robokassa/fail", func(c *gin.Context) {
		invIDStr := c.Query("InvId")

		c.JSON(200, gin.H{
			"ok":      false,
			"message": "Оплата отменена или произошла ошибка",
			"invId":   invIDStr,
		})
	})
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

func updateUserContentToPremium(ctx context.Context, userID, colorCode string, vacancies *repo.VacancyRepo, resumes *repo.ResumeRepo) error {
	// Обновляем все активные вакансии пользователя
	if err := vacancies.UpdateAllByCompanyID(ctx, userID, bson.M{
		"isPremium": true,
		"colorCode": colorCode,
	}); err != nil {
		return err
	}

	// Обновляем все активные резюме пользователя
	if err := resumes.UpdateAllByUserID(ctx, userID, bson.M{
		"isPremium": true,
		"colorCode": colorCode,
	}); err != nil {
		return err
	}

	return nil
}
