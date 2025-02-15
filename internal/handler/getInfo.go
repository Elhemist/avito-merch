package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getInfo(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, `no user id in token`)
		return
	}

	infoResp, err := h.services.User.GetInfo(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, infoResp)
}
