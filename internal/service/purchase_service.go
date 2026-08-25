package service

import (
	"fmt"
	"time"

	"github.com/purchase-decision-system/internal/domain"
	"github.com/purchase-decision-system/internal/infrastructure/aws"
	"github.com/purchase-decision-system/internal/infrastructure/config"
	"github.com/purchase-decision-system/internal/repository"
	log "github.com/sirupsen/logrus"
)

type PurchaseService struct {
	purchaseRepo *repository.PurchaseRepository
	decisionRepo *repository.DecisionRepository
	snsClient    *aws.SNSClient
	config       *config.BusinessConfig
}

func NewPurchaseService(
	purchaseRepo *repository.PurchaseRepository,
	decisionRepo *repository.DecisionRepository,
	snsClient *aws.SNSClient,
	config *config.BusinessConfig,
) *PurchaseService {
	return &PurchaseService{
		purchaseRepo: purchaseRepo,
		decisionRepo: decisionRepo,
		snsClient:    snsClient,
		config:       config,
	}
}

func (s *PurchaseService) ProcessPurchaseEvent(event *domain.PurchaseEvent) error {
	log.WithFields(log.Fields{
		"purchase_id": event.PurchaseID,
		"amount":      event.Amount,
		"customer_id": event.CustomerID,
	}).Info("Processing purchase event")

	if err := s.validatePurchaseEvent(event); err != nil {
		return fmt.Errorf("invalid purchase event: %w", err)
	}

	purchase := domain.NewPurchase(event)
	if err := s.purchaseRepo.Create(purchase); err != nil {
		return fmt.Errorf("failed to create purchase: %w", err)
	}

	decision := s.applyAutomaticRules(purchase)
	if decision != nil {
		if err := s.applyDecision(purchase, decision); err != nil {
			return fmt.Errorf("failed to apply automatic decision: %w", err)
		}
	}

	return nil
}

func (s *PurchaseService) ApprovePurchase(purchaseID, reason, approvedBy string) error {
	purchase, err := s.purchaseRepo.FindByID(purchaseID)
	if err != nil {
		return err
	}

	if purchase.Status != domain.StatusPending {
		return domain.ErrPurchaseAlreadyDecided
	}

	decision := domain.NewDecision(
		purchaseID,
		domain.StatusApproved,
		domain.DecisionManual,
		reason,
		approvedBy,
	)

	return s.applyDecision(purchase, decision)
}

func (s *PurchaseService) RejectPurchase(purchaseID, reason, rejectedBy string) error {
	purchase, err := s.purchaseRepo.FindByID(purchaseID)
	if err != nil {
		return err
	}

	if purchase.Status != domain.StatusPending {
		return domain.ErrPurchaseAlreadyDecided
	}

	decision := domain.NewDecision(
		purchaseID,
		domain.StatusRejected,
		domain.DecisionManual,
		reason,
		rejectedBy,
	)

	return s.applyDecision(purchase, decision)
}

func (s *PurchaseService) IgnorePurchase(purchaseID, reason, ignoredBy string) error {
	purchase, err := s.purchaseRepo.FindByID(purchaseID)
	if err != nil {
		return err
	}

	if purchase.Status != domain.StatusPending {
		return domain.ErrPurchaseAlreadyDecided
	}

	decision := domain.NewDecision(
		purchaseID,
		domain.StatusIgnored,
		domain.DecisionManual,
		reason,
		ignoredBy,
	)

	return s.applyDecision(purchase, decision)
}

func (s *PurchaseService) GetPurchase(purchaseID string) (*domain.Purchase, error) {
	return s.purchaseRepo.FindByID(purchaseID)
}

func (s *PurchaseService) GetPurchasesByStatus(status domain.PurchaseStatus, limit, offset int) ([]*domain.Purchase, error) {
	return s.purchaseRepo.FindByStatus(status, limit, offset)
}

func (s *PurchaseService) GetPurchaseHistory(purchaseID string) ([]*domain.Decision, error) {
	return s.decisionRepo.FindByPurchaseID(purchaseID)
}

func (s *PurchaseService) GetMetrics() (map[string]interface{}, error) {
	totalPurchases, _ := s.purchaseRepo.Count()
	pendingCount, _ := s.purchaseRepo.CountByStatus(domain.StatusPending)
	approvedCount, _ := s.purchaseRepo.CountByStatus(domain.StatusApproved)
	rejectedCount, _ := s.purchaseRepo.CountByStatus(domain.StatusRejected)
	ignoredCount, _ := s.purchaseRepo.CountByStatus(domain.StatusIgnored)

	manualDecisions, _ := s.decisionRepo.CountByDecisionType(domain.DecisionManual)
	automaticDecisions, _ := s.decisionRepo.CountByDecisionType(domain.DecisionAutomatic)

	metrics := map[string]interface{}{
		"total_purchases": totalPurchases,
		"by_status": map[string]int64{
			"pending":  pendingCount,
			"approved": approvedCount,
			"rejected": rejectedCount,
			"ignored":  ignoredCount,
		},
		"by_decision_type": map[string]int64{
			"manual":    manualDecisions,
			"automatic": automaticDecisions,
		},
		"approval_rate":  float64(approvedCount) / float64(totalPurchases) * 100,
		"rejection_rate": float64(rejectedCount) / float64(totalPurchases) * 100,
	}

	return metrics, nil
}

func (s *PurchaseService) validatePurchaseEvent(event *domain.PurchaseEvent) error {
	if event.Amount < s.config.MinPurchaseAmount {
		return domain.ErrInvalidAmount
	}
	if event.Amount > s.config.MaxPurchaseAmount {
		return domain.ErrInvalidAmount
	}
	if event.CustomerID == "" {
		return fmt.Errorf("customer_id is required")
	}
	if event.PurchaseID == "" {
		return fmt.Errorf("purchase_id is required")
	}
	return nil
}

func (s *PurchaseService) applyAutomaticRules(purchase *domain.Purchase) *domain.Decision {
	if purchase.Amount < s.config.AutoApproveThreshold {
		log.WithField("purchase_id", purchase.ID).Info("Auto-approving purchase (low amount)")
		return domain.NewDecision(
			purchase.ID,
			domain.StatusApproved,
			domain.DecisionAutomatic,
			fmt.Sprintf("Auto-approved: amount (%.2f) below threshold (%.2f)", purchase.Amount, s.config.AutoApproveThreshold),
			"system",
		)
	}

	if purchase.Amount > s.config.AutoRejectThreshold {
		log.WithField("purchase_id", purchase.ID).Info("Auto-rejecting purchase (high amount)")
		return domain.NewDecision(
			purchase.ID,
			domain.StatusRejected,
			domain.DecisionAutomatic,
			fmt.Sprintf("Auto-rejected: amount (%.2f) above threshold (%.2f)", purchase.Amount, s.config.AutoRejectThreshold),
			"system",
		)
	}

	log.WithField("purchase_id", purchase.ID).Info("Purchase requires manual review")
	return nil
}

func (s *PurchaseService) applyDecision(purchase *domain.Purchase, decision *domain.Decision) error {

	if err := s.decisionRepo.Create(decision); err != nil {
		return fmt.Errorf("failed to create decision: %w", err)
	}

	purchase.Status = decision.Status
	purchase.UpdatedAt = time.Now()
	if err := s.purchaseRepo.Update(purchase); err != nil {
		return fmt.Errorf("failed to update purchase: %w", err)
	}

	event := decision.ToDecisionEvent(purchase)
	if err := s.snsClient.PublishDecisionEvent(event); err != nil {
		log.WithError(err).Error("Failed to publish decision event to SNS")
	}

	log.WithFields(log.Fields{
		"purchase_id":   purchase.ID,
		"status":        decision.Status,
		"decision_type": decision.DecisionType,
		"decided_by":    decision.DecidedBy,
	}).Info("Decision applied successfully")

	return nil
}
