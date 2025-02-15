package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type coinTransfersRequest struct {
	ReceiverName string `json:"toUser"`
	Amount       int    `json:"amount"`
}

func (h *Handler) sendCoin(c *gin.Context) {
	var input coinTransfersRequest

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, `no user id in token`)
		return
	}
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, `binding error`)
		return
	}

	err = h.services.Wallet.SendCoin(userId, input.ReceiverName, input.Amount)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "ok")
}
