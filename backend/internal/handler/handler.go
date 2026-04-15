package handler

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	repo *repository.SubscriptionRepo
}

func NewSubscriptionHandler(repo *repository.SubscriptionRepo) *SubscriptionHandler {
	return &SubscriptionHandler{repo: repo}
}

// @Summary Create a new subscription
// @Description Create a new subscription record for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} domain.Subscription
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/subscriptions [post]
func (h *SubscriptionHandler) HandleCreateSubscription(c *gin.Context) {
	var req CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription := &domain.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate.Time,
	}

	if req.EndDate != nil && !req.EndDate.IsZero() {
		subscription.EndDate = &req.EndDate.Time
	}

	err := h.repo.CreateSubscription(subscription)
	if err != nil {
		slog.Error("failed to create subscription", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
		return
	}

	slog.Info("subscription created", "subscription_id", subscription.ID, "user_id", req.UserID)
	c.JSON(http.StatusCreated, subscription)
}

// @Summary Get a subscription by ID
// @Description Get a subscription record by its ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} domain.Subscription
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/subscriptions/{id} [get]
func (h *SubscriptionHandler) HandleGetSubscription(c *gin.Context) {
	var req GetSubscriptionRequest
	if err := c.ShouldBindUri(&req); err != nil {
		slog.Warn("invalid uri", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.repo.GetSubscriptionByID(req.ID)
	if err != nil {
		slog.Warn("subscription not found", "id", req.ID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	slog.Info("subscription retrieved", "id", req.ID)
	c.JSON(http.StatusOK, subscription)
}

// @Summary Update a subscription
// @Description Update an existing subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body UpdateSubscriptionBody false "Updated subscription data"
// @Success 200 {object} domain.Subscription
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/subscriptions/{id} [put]
func (h *SubscriptionHandler) HandleUpdateSubscription(c *gin.Context) {
	var req UpdateSubscriptionRequest
	if err := c.ShouldBindUri(&req); err != nil {
		slog.Warn("invalid uri", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.repo.GetSubscriptionByID(req.ID)
	if err != nil {
		slog.Warn("subscription not found", "id", req.ID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	var body UpdateSubscriptionBody
	if err := c.ShouldBindJSON(&body); err != nil && err.Error() != "EOF" {
		slog.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.ServiceName != nil {
		subscription.ServiceName = *body.ServiceName
	}
	if body.Price != nil {
		subscription.Price = *body.Price
	}
	if body.StartDate != nil && !body.StartDate.IsZero() {
		subscription.StartDate = body.StartDate.Time
	}
	if body.EndDate != nil {
		if body.EndDate.IsZero() {
			subscription.EndDate = nil
		} else {
			t := body.EndDate.Time
			subscription.EndDate = &t
		}
	}

	err = h.repo.UpdateSubscription(subscription)
	if err != nil {
		slog.Error("failed to update subscription", "id", req.ID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subscription"})
		return
	}

	slog.Info("subscription updated", "id", req.ID)
	c.JSON(http.StatusOK, subscription)
}

// @Summary Delete a subscription
// @Description Delete a subscription record by its ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/subscriptions/{id} [delete]
func (h *SubscriptionHandler) HandleDeleteSubscription(c *gin.Context) {
	var req DeleteSubscriptionRequest
	if err := c.ShouldBindUri(&req); err != nil {
		slog.Warn("invalid uri", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.repo.DeleteSubscription(req.ID)
	if err != nil {
		slog.Error("failed to delete subscription", "id", req.ID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete subscription"})
		return
	}

	slog.Info("subscription deleted", "id", req.ID)
	c.JSON(http.StatusNoContent, gin.H{"message": "subscription deleted successfully"})
}

// @Summary List all subscriptions
// @Description Get all subscriptions with optional user_id filter and pagination
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID filter"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} handler.PaginatedResponse{items=[]domain.Subscription}
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/subscriptions [get]
func (h *SubscriptionHandler) HandleGetSubscriptions(c *gin.Context) {
	var req GetSubscriptionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Warn("invalid query params", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page := req.Page
	if page < 1 {
		page = 1
	}

	limit := req.Limit
	if limit < 1 {
		limit = 10
	}

	subscriptions, total, err := h.repo.ListSubscriptions(req.UserID, page, limit)
	if err != nil {
		slog.Error("failed to list subscriptions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list subscriptions"})
		return
	}

	totalPages := (int(total) + limit - 1) / limit
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	response := PaginatedResponse{
		Items:      subscriptions,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	slog.Info("subscriptions listed", "count", len(subscriptions), "total", total, "page", page, "limit", limit, "user_id", req.UserID)
	c.JSON(http.StatusOK, response)
}

// @Summary Calculate total price
// @Description Calculate total cost of all subscriptions for a selected period with optional filters
// @Tags subscriptions
// @Produce json
// @Param period_start query string true "Period start (MM-YYYY)"
// @Param period_end query string true "Period end (MM-YYYY)"
// @Param user_id query string false "User ID filter"
// @Param service_name query string false "Service name filter"
// @Success 200 {object} gin.H{total_price=int}
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/calculate_total_price [get]
func (h *SubscriptionHandler) HandleCalculateTotalPrice(c *gin.Context) {
	var req CalculateTotalPriceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Warn("invalid query params", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	totalPrice, err := h.repo.CalculateTotalPrice(req.UserID, req.ServiceName, req.PeriodStart.Time, req.PeriodEnd.Time)
	if err != nil {
		slog.Error("failed to calculate total price", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate total price"})
		return
	}

	slog.Info("total price calculated", "total", totalPrice, "user_id", req.UserID, "service", req.ServiceName)
	c.JSON(http.StatusOK, gin.H{"total_price": totalPrice})
}
