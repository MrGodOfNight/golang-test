package handler

import (
	"net/http"
	"strconv"

	"testing-go/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Deposit(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.Deposit(c.Request.Context(), userID, req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "balance updated"})
}

func (h *UserHandler) Transfer(c *gin.Context) {
	var req struct {
		SenderID   int     `json:"sender_id"`
		ReceiverID int     `json:"receiver_id"`
		Amount     float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.Transfer(c.Request.Context(), req.SenderID, req.ReceiverID, req.Amount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transfer successful"})
}

func (h *UserHandler) GetTransactions(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("id"))
	transactions, err := h.service.GetLastTransactions(c.Request.Context(), userID, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
