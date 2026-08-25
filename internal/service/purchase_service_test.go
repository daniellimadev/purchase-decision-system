package service

import (
	"github.com/purchase-decision-system/internal/domain"
	"github.com/purchase-decision-system/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePurchaseEvent(t *testing.T) {
	service := &PurchaseService{
		config: &config.BusinessConfig{
			MinPurchaseAmount: 0.01,
			MaxPurchaseAmount: 10000.00,
		},
	}

	tests := []struct {
		name    string
		event   *domain.PurchaseEvent
		wantErr bool
	}{
		{
			name: "valid event",
			event: &domain.PurchaseEvent{
				PurchaseID: "test-123",
				CustomerID: "customer-123",
				Amount:     100.00,
				Currency:   "BRL",
			},
			wantErr: false,
		},
		{
			name: "amount too low",
			event: &domain.PurchaseEvent{
				PurchaseID: "test-123",
				CustomerID: "customer-123",
				Amount:     0.00,
			},
			wantErr: true,
		},
		{
			name: "amount too high",
			event: &domain.PurchaseEvent{
				PurchaseID: "test-123",
				CustomerID: "customer-123",
				Amount:     20000.00,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validatePurchaseEvent(tt.event)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
