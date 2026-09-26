package handlers

import (
	"astroolog/backend/internal/repositories"
	"astroolog/backend/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Handler struct {
	Repo     repositories.Repository
	Requests services.Requests
}

func fail(c *gin.Context, err error) {
	log.Printf("API error: %v", err)
	c.JSON(500, gin.H{"error": "Не удалось обработать запрос. Попробуйте позже.", "code": "server_error"})
}
func (h Handler) ListServices(c *gin.Context) {
	v, err := h.Repo.Services()
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, v)
}
func (h Handler) Service(c *gin.Context) {
	v, err := h.Repo.Service(c.Param("slug"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{"error": "Услуга не найдена", "code": "not_found"})
		return
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, v)
}
func (h Handler) Submit(c *gin.Context, booking bool) {
	var input services.Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Проверьте имя, email, длину полей и согласие на обработку данных.", "code": "validation_error"})
		return
	}
	err := h.Requests.Submit(input, booking, c.GetHeader("Idempotency-Key"))
	if errors.Is(err, services.ErrStalePrice) {
		c.JSON(409, gin.H{"code": "stale_price"})
		return
	}
	if errors.Is(err, services.ErrIdempotency) {
		c.JSON(409, gin.H{"code": "idempotency_conflict"})
		return
	}

	if errors.Is(err, services.ErrInvalid) {
		c.JSON(400, gin.H{"error": "Проверьте поля: услуга должна быть доступна, дата — не в прошлом, сообщение — заполнено.", "code": "invalid_request"})
		return
	}
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Спасибо! Заявка сохранена. Мы свяжемся с вами по указанному email.", "code": "request_created"})
}

func (h Handler) BookingOptions(c *gin.Context) {
	options, err := h.Repo.BookingOptions()
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(200, options)
}
