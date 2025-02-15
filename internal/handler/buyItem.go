package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) buyItem(c *gin.Context) {
	itemStr := c.Param("item")
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, `no user id in token`)
		return
	}
	err = h.services.Inventory.BuyItem(userId, itemStr)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "ok")
}
