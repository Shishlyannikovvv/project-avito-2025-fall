package handler

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/domain"
	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/service"
	"github.com/go-chi/chi/v5"
)

// Handler управляет HTTP-запросами
type Handler struct {
	srv *service.Service
}

// NewHandler создаёт новый хэндлер
func NewHandler(srv *service.Service) *Handler {
	return &Handler{srv: srv}
}

// SetupRoutes настраивает роуты
func (h *Handler) SetupRoutes(r *chi.Mux) {
	r.Post("/teams", h.createTeam)
	r.Post("/users", h.createUser)
	r.Post("/pull-requests", h.createPR)
	r.Patch("/pull-requests/{prId}/reassign", h.reassignReviewer)
	r.Patch("/pull-requests/{prId}/merge", h.mergePR)
	r.Get("/users/{userId}/pull-requests", h.getUserPRs)

	// Бонус: статистика
	r.Get("/stats/reviewers", h.getReviewerStats)

	// Бонус: массовая деактивация
	r.Patch("/teams/{teamId}/deactivate", h.deactivateTeam)
}

func (h *Handler) createTeam(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	team, err := h.srv.CreateTeam(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(team)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	user, err := h.srv.CreateUser(r.Context(), req.Username, req.TeamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) createPR(w http.ResponseWriter, r *http.Request) {
	var req domain.CreatePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	pr, err := h.srv.CreatePR(r.Context(), req.Title, req.AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pr)
}

func (h *Handler) reassignReviewer(w http.ResponseWriter, r *http.Request) {
	prID, _ := strconv.Atoi(chi.URLParam(r, "prId"))
	var req domain.ReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	pr, err := h.srv.ReassignReviewer(r.Context(), prID, req.ReviewerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pr)
}

func (h *Handler) mergePR(w http.ResponseWriter, r *http.Request) {
	prID, _ := strconv.Atoi(chi.URLParam(r, "prId"))
	pr, err := h.srv.MergePR(r.Context(), prID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pr)
}

func (h *Handler) getUserPRs(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(chi.URLParam(r, "userId"))
	prs, err := h.srv.GetUserPRs(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(prs)
}

// Бонус: статистика по ревьюверам
func (h *Handler) getReviewerStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.srv.GetReviewerStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}

// Бонус: деактивация команды
func (h *Handler) deactivateTeam(w http.ResponseWriter, r *http.Request) {
	teamID, _ := strconv.Atoi(chi.URLParam(r, "teamId"))
	err := h.srv.DeactivateTeam(r.Context(), teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}