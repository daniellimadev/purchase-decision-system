package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/purchase-decision-system/internal/domain"
	"github.com/purchase-decision-system/internal/services"
	log "github.com/sirupsen/logrus"
)

type PurchaseHandler struct {
	service *services.PurchaseService
}

func NewPurchaseHandler(service *services.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{service: service}
}

type ApproveRequest struct {
	Reason     string `json:"reason" binding:"required"`
	ApprovedBy string `json:"approved_by" binding:"required"`
}

type RejectRequest struct {
	Reason     string `json:"reason" binding:"required"`
	RejectedBy string `json:"rejected_by" binding:"required"`
}

type IgnoreRequest struct {
	Reason    string `json:"reason" binding:"required"`
	IgnoredBy string `json:"ignored_by" binding:"required"`
}

func (h *PurchaseHandler) ApprovePurchase(c *gin.Context) {
	purchaseID := c.Param("id")

	var req ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ApprovePurchase(purchaseID, req.Reason, req.ApprovedBy); err != nil {
		log.WithError(err).Error("Failed to approve purchase")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Purchase approved successfully",
		"purchase_id": purchaseID,
		"status":      domain.StatusApproved,
	})
}

func (h *PurchaseHandler) RejectPurchase(c *gin.Context) {
	purchaseID := c.Param("id")

	var req RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RejectPurchase(purchaseID, req.Reason, req.RejectedBy); err != nil {
		log.WithError(err).Error("Failed to reject purchase")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Purchase rejected successfully",
		"purchase_id": purchaseID,
		"status":      domain.StatusRejected,
	})
}

func (h *PurchaseHandler) IgnorePurchase(c *gin.Context) {
	purchaseID := c.Param("id")

	var req IgnoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.IgnorePurchase(purchaseID, req.Reason, req.IgnoredBy); err != nil {
		log.WithError(err).Error("Failed to ignore purchase")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Purchase ignored successfully",
		"purchase_id": purchaseID,
		"status":      domain.StatusIgnored,
	})
}

func (h *PurchaseHandler) GetPurchase(c *gin.Context) {
	purchaseID := c.Param("id")

	purchase, err := h.service.GetPurchase(purchaseID)
	if err != nil {
		if err == domain.ErrPurchaseNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Purchase not found"})
			return
		}
		log.WithError(err).Error("Failed to get purchase")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

func (h *PurchaseHandler) ListPurchases(c *gin.Context) {
	status := c.DefaultQuery("status", string(domain.StatusPending))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	purchases, err := h.service.GetPurchasesByStatus(domain.PurchaseStatus(status), limit, offset)
	if err != nil {
		log.WithError(err).Error("Failed to list purchases")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"purchases": purchases,
		"count":     len(purchases),
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *PurchaseHandler) GetPurchaseHistory(c *gin.Context) {
	purchaseID := c.Param("id")

	history, err := h.service.GetPurchaseHistory(purchaseID)
	if err != nil {
		log.WithError(err).Error("Failed to get purchase history")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"purchase_id": purchaseID,
		"history":     history,
		"count":       len(history),
	})
}

func (h *PurchaseHandler) GetMetrics(c *gin.Context) {
	metrics, err := h.service.GetMetrics()
	if err != nil {
		log.WithError(err).Error("Failed to get metrics")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

func (h *PurchaseHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "purchase-decision-system",
	})
}
