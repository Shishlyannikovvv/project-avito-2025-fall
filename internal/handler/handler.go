package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Shishlyannikovvv/project-avito-2025-fall/internal/service"
	"github.com/go-chi/chi/v5"
)

// HandlerImpl — наша реализация интерфейса из generated.go
type HandlerImpl struct {
	srv *service.Service
}

func NewHandler(srv *service.Service) *HandlerImpl {
	return &HandlerImpl{srv: srv}
}

// Реализация методов из StrictServerInterface (имена точно такие же, как в generated.go)

func (h *HandlerImpl) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	team, err := h.srv.CreateTeam(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(team)
}

func (h *HandlerImpl) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		TeamID   int    `json:"team_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	user, err := h.srv.CreateUser(r.Context(), req.Username, req.TeamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *HandlerImpl) CreatePullRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string `json:"title"`
		AuthorID int    `json:"author_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	pr, err := h.srv.CreatePR(r.Context(), req.Title, req.AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pr)
}

func (h *HandlerImpl) ReassignReviewer(w http.ResponseWriter, r *http.Request, prId int) {
	var req struct {
		ReviewerID int `json:"reviewer_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	pr, err := h.srv.ReassignReviewer(r.Context(), prId, req.ReviewerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pr)
}

func (h *HandlerImpl) MergePullRequest(w http.ResponseWriter, r *http.Request, prId int) {
	pr, err := h.srv.MergePR(r.Context(), prId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(pr)
}

func (h *HandlerImpl) GetUserPullRequests(w http.ResponseWriter, r *http.Request, userId int) {
	prs, err := h.srv.GetUserPRs(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(prs)
}

// Бонусные — можно оставить, они не ломаются
func (h *HandlerImpl) GetReviewerStats(w http.ResponseWriter, r *http.Request) {
	stats, _ := h.srv.GetReviewerStats(r.Context())
	json.NewEncoder(w).Encode(stats)
}

func (h *HandlerImpl) DeactivateTeam(w http.ResponseWriter, r *http.Request) {
	teamId, _ := strconv.Atoi(chi.URLParam(r, "teamId"))
	_ = h.srv.DeactivateTeam(r.Context(), teamId)
	w.WriteHeader(http.StatusOK)
}
