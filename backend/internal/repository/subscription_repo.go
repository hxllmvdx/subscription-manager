package repository

import (
	"backend/internal/domain"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	CreateSubscription(subscription *domain.Subscription) error
	GetSubscriptionByID(id string) (*domain.Subscription, error)
	UpdateSubscription(subscription *domain.Subscription) error
	DeleteSubscription(id string) error
	ListSubscriptions(userID string, page, limit int) ([]domain.Subscription, int64, error)
	CalculateTotalPrice(userID, serviceName string, periodStart, periodEnd time.Time) (int, error)
}

type SubscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepo(db *gorm.DB) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) CreateSubscription(subscription *domain.Subscription) error {
	slog.Debug("creating subscription", "user_id", subscription.UserID, "service", subscription.ServiceName)
	err := r.db.Create(subscription).Error
	if err != nil {
		slog.Error("failed to create subscription in db", "error", err)
		return err
	}
	return nil
}

func (r *SubscriptionRepo) GetSubscriptionByID(id string) (*domain.Subscription, error) {
	slog.Debug("getting subscription", "id", id)
	var subscription domain.Subscription
	err := r.db.First(&subscription, "id = ?", id).Error
	if err != nil {
		slog.Warn("subscription not found in db", "id", id, "error", err)
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepo) UpdateSubscription(subscription *domain.Subscription) error {
	slog.Debug("updating subscription", "id", subscription.ID)
	err := r.db.Save(subscription).Error
	if err != nil {
		slog.Error("failed to update subscription in db", "id", subscription.ID, "error", err)
		return err
	}
	return nil
}

func (r *SubscriptionRepo) DeleteSubscription(id string) error {
	slog.Debug("deleting subscription", "id", id)
	err := r.db.Delete(&domain.Subscription{}, "id = ?", id).Error
	if err != nil {
		slog.Error("failed to delete subscription from db", "id", id, "error", err)
		return err
	}
	return nil
}

func (r *SubscriptionRepo) ListSubscriptions(userID string, page, limit int) ([]domain.Subscription, int64, error) {
	var subscriptions []domain.Subscription
	query := r.db.Model(&domain.Subscription{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		slog.Error("failed to count subscriptions", "error", err)
		return nil, 0, err
	}

	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Order("start_date DESC").Find(&subscriptions).Error
	if err != nil {
		slog.Error("failed to list subscriptions", "error", err)
		return nil, 0, err
	}
	return subscriptions, total, nil
}

func (r *SubscriptionRepo) CalculateTotalPrice(userID, serviceName string, periodStart, periodEnd time.Time) (int, error) {
	slog.Debug("calculating total price", "user_id", userID, "service", serviceName, "period_start", periodStart, "period_end", periodEnd)

	query := r.db.Model(&domain.Subscription{})
	query = query.Where("start_date <= ?", periodEnd)
	query = query.Where("(end_date >= ? OR end_date IS NULL)", periodStart)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	var subscriptions []domain.Subscription
	err := query.Find(&subscriptions).Error
	if err != nil {
		slog.Error("failed to fetch subscriptions for price calculation", "error", err)
		return 0, err
	}

	var totalPrice int
	for _, sub := range subscriptions {
		actualStart := sub.StartDate
		actualEnd := sub.EndDate
		if actualEnd == nil || actualEnd.After(periodEnd) {
			actualEnd = &periodEnd
		}

		overlapStart := actualStart
		if periodStart.After(overlapStart) {
			overlapStart = periodStart
		}

		overlapEnd := *actualEnd
		if periodEnd.Before(overlapEnd) {
			overlapEnd = periodEnd
		}

		months := calculateMonthsOverlap(overlapStart, overlapEnd)
		if months <= 0 {
			continue
		}

		totalPrice += sub.Price * months
	}

	slog.Debug("total price calculated", "total", totalPrice, "subscriptions_count", len(subscriptions))
	return totalPrice, nil
}

func calculateMonthsOverlap(start, end time.Time) int {
	if end.Before(start) {
		return 0
	}

	months := (end.Year()-start.Year())*12 + int(end.Month()) - int(start.Month())

	if end.Day() < start.Day() {
		months--
	}

	if months < 0 {
		return 0
	}

	if months == 0 {
		return 1
	}

	return months
}
