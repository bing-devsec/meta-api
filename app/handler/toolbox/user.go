package toolbox

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"meta-api/common/codes"
	"meta-api/common/types"
)

func (h *toolboxHandler) UserGetToolboxList(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	ctx := c.Request.Context()
	response, err := h.service.UserGetToolboxList(ctx)
	if err != nil {
		c.JSON(http.StatusOK, types.Response{Code: codes.InternalServerError, Message: "获取工具箱失败", Data: nil})
		return
	}
	c.JSON(http.StatusOK, types.Response{Code: codes.Success, Message: "", Data: response})
}
