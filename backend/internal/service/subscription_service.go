package service

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"errors"
	"log/slog"
	"time"
)

type SubscriptionService interface {
	CreateSubscription(req *CreateSubscriptionRequest) (*domain.Subscription, error)
	GetSubscriptionByID(id string) (*domain.Subscription, error)
	UpdateSubscription(id string, req *UpdateSubscriptionRequest) (*domain.Subscription, error)
	DeleteSubscription(id string) error
	ListSubscriptions(req *ListSubscriptionsRequest) ([]domain.Subscription, int64, int, error)
	CalculateTotalPrice(req *CalculateTotalPriceRequest) (int, error)
}

type CreateSubscriptionRequest struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
	EndDate     *time.Time
}

type UpdateSubscriptionRequest struct {
	ServiceName *string
	Price       *int
	StartDate   *time.Time
	EndDate     *time.Time
}

type ListSubscriptionsRequest struct {
	UserID string
	Page   int
	Limit  int
}

type CalculateTotalPriceRequest struct {
	UserID      string
	ServiceName string
	PeriodStart time.Time
	PeriodEnd   time.Time
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(req *CreateSubscriptionRequest) (*domain.Subscription, error) {
	if req.ServiceName == "" {
		return nil, errors.New("service_name is required")
	}
	if req.Price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}
	if req.UserID == "" {
		return nil, errors.New("user_id is required")
	}
	if req.StartDate.IsZero() {
		return nil, errors.New("start_date is required")
	}

	subscription := &domain.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	err := s.repo.CreateSubscription(subscription)
	if err != nil {
		slog.Error("failed to create subscription", "error", err)
		return nil, err
	}

	slog.Info("subscription created", "subscription_id", subscription.ID, "user_id", req.UserID)
	return subscription, nil
}

func (s *subscriptionService) GetSubscriptionByID(id string) (*domain.Subscription, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	subscription, err := s.repo.GetSubscriptionByID(id)
	if err != nil {
		slog.Warn("subscription not found", "id", id, "error", err)
		return nil, err
	}

	slog.Info("subscription retrieved", "id", id)
	return subscription, nil
}

func (s *subscriptionService) UpdateSubscription(id string, req *UpdateSubscriptionRequest) (*domain.Subscription, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	subscription, err := s.repo.GetSubscriptionByID(id)
	if err != nil {
		slog.Warn("subscription not found", "id", id, "error", err)
		return nil, err
	}

	if req.ServiceName != nil {
		subscription.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, errors.New("price must be greater than 0")
		}
		subscription.Price = *req.Price
	}
	if req.StartDate != nil && !req.StartDate.IsZero() {
		subscription.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		if req.EndDate.IsZero() {
			subscription.EndDate = nil
		} else {
			t := *req.EndDate
			subscription.EndDate = &t
		}
	}

	err = s.repo.UpdateSubscription(subscription)
	if err != nil {
		slog.Error("failed to update subscription", "id", id, "error", err)
		return nil, err
	}

	slog.Info("subscription updated", "id", id)
	return subscription, nil
}

func (s *subscriptionService) DeleteSubscription(id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	err := s.repo.DeleteSubscription(id)
	if err != nil {
		slog.Error("failed to delete subscription", "id", id, "error", err)
		return err
	}

	slog.Info("subscription deleted", "id", id)
	return nil
}

func (s *subscriptionService) ListSubscriptions(req *ListSubscriptionsRequest) ([]domain.Subscription, int64, int, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}

	limit := req.Limit
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	subscriptions, total, err := s.repo.ListSubscriptions(req.UserID, page, limit)
	if err != nil {
		slog.Error("failed to list subscriptions", "error", err)
		return nil, 0, 0, err
	}

	totalPages := (int(total) + limit - 1) / limit
	if totalPages == 0 && total > 0 {
		totalPages = 1
	}

	slog.Info("subscriptions listed", "count", len(subscriptions), "total", total, "page", page, "limit", limit, "user_id", req.UserID)
	return subscriptions, total, totalPages, nil
}

func (s *subscriptionService) CalculateTotalPrice(req *CalculateTotalPriceRequest) (int, error) {
	if req.PeriodStart.IsZero() || req.PeriodEnd.IsZero() {
		return 0, errors.New("period_start and period_end are required")
	}

	totalPrice, err := s.repo.CalculateTotalPrice(req.UserID, req.ServiceName, req.PeriodStart, req.PeriodEnd)
	if err != nil {
		slog.Error("failed to calculate total price", "error", err)
		return 0, err
	}

	slog.Info("total price calculated", "total", totalPrice, "user_id", req.UserID, "service", req.ServiceName)
	return totalPrice, nil
}
