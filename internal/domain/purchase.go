package domain

import (
	"github.com/google/uuid"
	"time"
)

type PurchaseStatus string

const (
	StatusPending  PurchaseStatus = "pending"
	StatusApproved PurchaseStatus = "approved"
	StatusRejected PurchaseStatus = "rejected"
	StatusIgnored  PurchaseStatus = "ignored"
)

type DecisionType string

const (
	DecisionManual    DecisionType = "manual"
	DecisionAutomatic DecisionType = "automatic"
)

type Purchase struct {
	ID            string                 `json:"id" gorm:"primaryKey"`
	CustomerID    string                 `json:"customer_id" gorm:"index"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	Merchant      string                 `json:"merchant"`
	PaymentMethod string                 `json:"payment_method"`
	Status        PurchaseStatus         `json:"status" gorm:"index"`
	Metadata      map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type Decision struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	PurchaseID   string         `json:"purchase_id" gorm:"index"`
	Status       PurchaseStatus `json:"status"`
	DecisionType DecisionType   `json:"decision_type"`
	Reason       string         `json:"reason"`
	DecidedBy    string         `json:"decided_by"`
	CreatedAt    time.Time      `json:"created_at"`
	Purchase     *Purchase      `json:"purchase,omitempty" gorm:"foreignKey:PurchaseID"`
}

type PurchaseEvent struct {
	PurchaseID    string                 `json:"purchase_id"`
	CustomerID    string                 `json:"customer_id"`
	Amount        float64                `json:"amount"`
	Currency      string                 `json:"currency"`
	Merchant      string                 `json:"merchant"`
	PaymentMethod string                 `json:"payment_method"`
	Timestamp     time.Time              `json:"timestamp"`
	Metadata      map[string]interface{} `json:"metadata"`
}

type DecisionEvent struct {
	EventID      string                 `json:"event_id"`
	PurchaseID   string                 `json:"purchase_id"`
	CustomerID   string                 `json:"customer_id"`
	Status       PurchaseStatus         `json:"status"`
	DecisionType DecisionType           `json:"decision_type"`
	Reason       string                 `json:"reason"`
	DecidedBy    string                 `json:"decided_by"`
	Amount       float64                `json:"amount"`
	Timestamp    time.Time              `json:"timestamp"`
	Metadata     map[string]interface{} `json:"metadata"`
}

func NewPurchase(event *PurchaseEvent) *Purchase {
	return &Purchase{
		ID:            event.PurchaseID,
		CustomerID:    event.CustomerID,
		Amount:        event.Amount,
		Currency:      event.Currency,
		Merchant:      event.Merchant,
		PaymentMethod: event.PaymentMethod,
		Status:        StatusPending,
		Metadata:      event.Metadata,
		CreatedAt:     event.Timestamp,
		UpdatedAt:     event.Timestamp,
	}
}

func NewDecision(purchaseID string, status PurchaseStatus, decisionType DecisionType, reason, decidedBy string) *Decision {
	return &Decision{
		ID:           uuid.New().String(),
		PurchaseID:   purchaseID,
		Status:       status,
		DecisionType: decisionType,
		Reason:       reason,
		DecidedBy:    decidedBy,
		CreatedAt:    time.Now(),
	}
}

func (d *Decision) ToDecisionEvent(purchase *Purchase) *DecisionEvent {
	return &DecisionEvent{
		EventID:      uuid.New().String(),
		PurchaseID:   d.PurchaseID,
		CustomerID:   purchase.CustomerID,
		Status:       d.Status,
		DecisionType: d.DecisionType,
		Reason:       d.Reason,
		DecidedBy:    d.DecidedBy,
		Amount:       purchase.Amount,
		Timestamp:    d.CreatedAt,
		Metadata:     purchase.Metadata,
	}
}
