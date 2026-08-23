package repository

import (
	"fmt"

	"github.com/purchase-decision-system/internal/domain"
	"gorm.io/gorm"
)

type PurchaseRepository struct {
	db *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) Create(purchase *domain.Purchase) error {
	if err := r.db.Create(purchase).Error; err != nil {
		return fmt.Errorf("failed to create purchase: %w", err)
	}
	return nil
}

func (r *PurchaseRepository) FindByID(id string) (*domain.Purchase, error) {
	var purchase domain.Purchase
	if err := r.db.First(&purchase, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrPurchaseNotFound
		}
		return nil, fmt.Errorf("failed to find purchase: %w", err)
	}
	return &purchase, nil
}

func (r *PurchaseRepository) FindByStatus(status domain.PurchaseStatus, limit, offset int) ([]*domain.Purchase, error) {
	var purchases []*domain.Purchase
	if err := r.db.Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&purchases).Error; err != nil {
		return nil, fmt.Errorf("failed to find purchases: %w", err)
	}
	return purchases, nil
}

func (r *PurchaseRepository) FindByCustomerID(customerID string, limit, offset int) ([]*domain.Purchase, error) {
	var purchases []*domain.Purchase
	if err := r.db.Where("customer_id = ?", customerID).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&purchases).Error; err != nil {
		return nil, fmt.Errorf("failed to find purchases: %w", err)
	}
	return purchases, nil
}

func (r *PurchaseRepository) Update(purchase *domain.Purchase) error {
	if err := r.db.Save(purchase).Error; err != nil {
		return fmt.Errorf("failed to update purchase: %w", err)
	}
	return nil
}

func (r *PurchaseRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&domain.Purchase{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count purchases: %w", err)
	}
	return count, nil
}

func (r *PurchaseRepository) CountByStatus(status domain.PurchaseStatus) (int64, error) {
	var count int64
	if err := r.db.Model(&domain.Purchase{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count purchases: %w", err)
	}
	return count, nil
}
