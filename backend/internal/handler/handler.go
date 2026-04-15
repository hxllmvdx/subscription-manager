package handler

import (
	"backend/internal/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	svc service.SubscriptionService
}

func NewSubscriptionHandler(svc service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc}
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

	serviceReq := &service.CreateSubscriptionRequest{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate.Time,
	}

	if req.EndDate != nil && !req.EndDate.IsZero() {
		t := req.EndDate.Time
		serviceReq.EndDate = &t
	}

	subscription, err := h.svc.CreateSubscription(serviceReq)
	if err != nil {
		slog.Error("failed to create subscription", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
		return
	}

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

	subscription, err := h.svc.GetSubscriptionByID(req.ID)
	if err != nil {
		slog.Warn("subscription not found", "id", req.ID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

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
	var uriReq UpdateSubscriptionRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		slog.Warn("invalid uri", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body UpdateSubscriptionBody
	if err := c.ShouldBindJSON(&body); err != nil && err.Error() != "EOF" {
		slog.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serviceReq := &service.UpdateSubscriptionRequest{}

	if body.ServiceName != nil {
		serviceReq.ServiceName = body.ServiceName
	}
	if body.Price != nil {
		serviceReq.Price = body.Price
	}
	if body.StartDate != nil && !body.StartDate.IsZero() {
		t := body.StartDate.Time
		serviceReq.StartDate = &t
	}
	if body.EndDate != nil {
		if body.EndDate.IsZero() {
			serviceReq.EndDate = nil
		} else {
			t := body.EndDate.Time
			serviceReq.EndDate = &t
		}
	}

	subscription, err := h.svc.UpdateSubscription(uriReq.ID, serviceReq)
	if err != nil {
		slog.Error("failed to update subscription", "id", uriReq.ID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subscription"})
		return
	}

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

	err := h.svc.DeleteSubscription(req.ID)
	if err != nil {
		slog.Error("failed to delete subscription", "id", req.ID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete subscription"})
		return
	}

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

	serviceReq := &service.ListSubscriptionsRequest{
		UserID: req.UserID,
		Page:   req.Page,
		Limit:  req.Limit,
	}

	subscriptions, total, totalPages, err := h.svc.ListSubscriptions(serviceReq)
	if err != nil {
		slog.Error("failed to list subscriptions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list subscriptions"})
		return
	}

	response := PaginatedResponse{
		Items:      subscriptions,
		Total:      total,
		Page:       serviceReq.Page,
		Limit:      serviceReq.Limit,
		TotalPages: totalPages,
	}

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

	serviceReq := &service.CalculateTotalPriceRequest{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		PeriodStart: req.PeriodStart.Time,
		PeriodEnd:   req.PeriodEnd.Time,
	}

	totalPrice, err := h.svc.CalculateTotalPrice(serviceReq)
	if err != nil {
		slog.Error("failed to calculate total price", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate total price"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_price": totalPrice})
}
