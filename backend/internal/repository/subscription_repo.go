package repository

import (
	"backend/internal/domain"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

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

func (r *SubscriptionRepo) ListSubscriptions(userID string) ([]domain.Subscription, error) {
	var subscriptions []domain.Subscription
	query := r.db.Model(&domain.Subscription{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Order("start_date DESC").Find(&subscriptions).Error
	if err != nil {
		slog.Error("failed to list subscriptions", "error", err)
		return nil, err
	}
	return subscriptions, nil
}

func (r *SubscriptionRepo) CalculateTotalPrice(userID, serviceName string, periodStart, periodEnd time.Time) (int, error) {
	slog.Debug("calculating total price", "user_id", userID, "service", serviceName, "period_start", periodStart, "period_end", periodEnd)

	query := r.db.Model(&domain.Subscription{}).Select("COALESCE(SUM(price), 0)")

	query = query.Where("start_date <= ?", periodEnd)
	query = query.Where("(end_date >= ? OR end_date IS NULL)", periodStart)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	var totalPrice int
	err := query.Scan(&totalPrice).Error
	if err != nil {
		slog.Error("failed to calculate total price", "error", err)
		return 0, err
	}

	return totalPrice, nil
}
