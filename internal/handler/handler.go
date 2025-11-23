package handler

import (
	"net/http"

	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/domain"
	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// Вспомогательная структура для парсинга JSON при создании PR
type CreatePRRequest struct {
	Title    string `json:"title" binding:"required"`
	AuthorID int    `json:"author_id" binding:"required"`
}

func (h *Handler) CreateTeam(c *gin.Context) {
	var req domain.Team
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	team, err := h.svc.CreateTeam(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, team)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req domain.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.CreateUser(req.Name, req.TeamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) CreatePR(c *gin.Context) {
	var req CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	pr, err := h.svc.CreatePR(req.Title, req.AuthorID)
	if err != nil {
		// В реальном проекте тут нужно различать ошибки (404 vs 500)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pr)
}

// RegisterRoutes регистрирует пути в Gin
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/teams", h.CreateTeam)
		api.POST("/users", h.CreateUser)
		api.POST("/pull-requests", h.CreatePR)
	}
}
