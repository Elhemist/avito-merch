package handler

import (
	"merch-test/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) authorization(c *gin.Context) {
	var input models.AuthRequest

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, `binding error`)
		return
	}

	token, err := h.services.Authorization.GenerateToken(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, `password check error`)
		return
	}

	logrus.Info("user logged: ", input.Username)
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
