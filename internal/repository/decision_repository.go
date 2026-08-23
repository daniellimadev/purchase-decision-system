package repository

import (
	"fmt"

	"github.com/purchase-decision-system/internal/domain"
	"gorm.io/gorm"
)

type DecisionRepository struct {
	db *gorm.DB
}

func NewDecisionRepository(db *gorm.DB) *DecisionRepository {
	return &DecisionRepository{db: db}
}

func (r *DecisionRepository) Create(decision *domain.Decision) error {
	if err := r.db.Create(decision).Error; err != nil {
		return fmt.Errorf("failed to create decision: %w", err)
	}
	return nil
}

func (r *DecisionRepository) FindByPurchaseID(purchaseID string) ([]*domain.Decision, error) {
	var decisions []*domain.Decision
	if err := r.db.Where("purchase_id = ?", purchaseID).
		Order("created_at DESC").
		Find(&decisions).Error; err != nil {
		return nil, fmt.Errorf("failed to find decisions: %w", err)
	}
	return decisions, nil
}

func (r *DecisionRepository) FindLatestByPurchaseID(purchaseID string) (*domain.Decision, error) {
	var decision domain.Decision
	if err := r.db.Where("purchase_id = ?", purchaseID).
		Order("created_at DESC").
		First(&decision).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find decision: %w", err)
	}
	return &decision, nil
}

func (r *DecisionRepository) FindByStatus(status domain.PurchaseStatus, limit, offset int) ([]*domain.Decision, error) {
	var decisions []*domain.Decision
	if err := r.db.Preload("Purchase").
		Where("status = ?", status).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&decisions).Error; err != nil {
		return nil, fmt.Errorf("failed to find decisions: %w", err)
	}
	return decisions, nil
}

func (r *DecisionRepository) CountByStatus(status domain.PurchaseStatus) (int64, error) {
	var count int64
	if err := r.db.Model(&domain.Decision{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count decisions: %w", err)
	}
	return count, nil
}

func (r *DecisionRepository) CountByDecisionType(decisionType domain.DecisionType) (int64, error) {
	var count int64
	if err := r.db.Model(&domain.Decision{}).Where("decision_type = ?", decisionType).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count decisions: %w", err)
	}
	return count, nil
}
