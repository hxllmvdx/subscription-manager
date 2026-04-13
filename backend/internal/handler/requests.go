package handler

import (
	"fmt"
	"strings"
	"time"
)

// MonthYearDate represents a date in "MM-YYYY" format
type MonthYearDate struct {
	time.Time
}

func (m *MonthYearDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" || s == "null" {
		return nil
	}

	t, err := time.Parse("01-2006", s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected MM-YYYY: %w", err)
	}

	m.Time = t
	return nil
}

func (m *MonthYearDate) MarshalJSON() ([]byte, error) {
	if m.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, m.Time.Format("01-2006"))), nil
}

type CreateSubscriptionRequest struct {
	ServiceName string        `json:"service_name" binding:"required"`
	Price       int           `json:"price" binding:"required,gt=0"`
	UserID      string        `json:"user_id" binding:"required,uuid"`
	StartDate   MonthYearDate `json:"start_date" binding:"required"`
	EndDate     *MonthYearDate `json:"end_date"`
}

type GetSubscriptionRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type UpdateSubscriptionRequest struct {
	ID          string         `uri:"id" binding:"required,uuid"`
	ServiceName string         `json:"service_name"`
	Price       *int           `json:"price" binding:"omitempty,gt=0"`
	StartDate   *MonthYearDate `json:"start_date"`
	EndDate     *MonthYearDate `json:"end_date"`
}

type DeleteSubscriptionRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type UpdateSubscriptionBody struct {
	ServiceName *string         `json:"service_name"`
	Price       *int            `json:"price" binding:"omitempty,gt=0"`
	StartDate   *MonthYearDate  `json:"start_date"`
	EndDate     *MonthYearDate  `json:"end_date"`
}

type GetSubscriptionsRequest struct {
	UserID string `form:"user_id" binding:"omitempty,uuid"`
}

type CalculateTotalPriceRequest struct {
	UserID      string        `form:"user_id" binding:"omitempty,uuid"`
	ServiceName string        `form:"service_name"`
	PeriodStart MonthYearDate `form:"period_start" binding:"required"`
	PeriodEnd   MonthYearDate `form:"period_end" binding:"required"`
}
