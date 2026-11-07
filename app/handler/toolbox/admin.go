package toolbox

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"meta-api/common/codes"
	"meta-api/common/idutil"
	"meta-api/common/types"
)

const toolboxURLMaxBytes = 2048

func (h *toolboxHandler) AdminGetToolboxList(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	ctx := c.Request.Context()
	response, err := h.service.AdminGetToolboxList(ctx)
	if err != nil {
		c.JSON(http.StatusOK, types.Response{Code: codes.InternalServerError, Message: "获取工具箱失败", Data: nil})
		return
	}
	c.JSON(http.StatusOK, types.Response{Code: codes.Success, Message: "", Data: response})
}

func (h *toolboxHandler) AdminAddToolboxTool(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	ctx := c.Request.Context()
	request := new(types.AdminAddToolboxToolRequest)
	if err := c.ShouldBindJSON(request); err != nil {
		h.logger.Warn("parameter binding error", zap.Error(err))
		message := "无效的请求参数"
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
			field := validationErrors[0]
			if field.Field() == "Name" && field.Tag() == "max" {
				message = "工具名称应为 1 至 30 个字符"
			} else if field.Field() == "Description" && field.Tag() == "max" {
				message = "工具简介最多 100 个字符"
			}
		}
		c.JSON(http.StatusOK, types.Response{Code: codes.BadRequest, Message: message, Data: nil})
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" {
		c.JSON(http.StatusOK, types.Response{Code: codes.BadRequest, Message: "工具名称应为 1 至 30 个字符", Data: nil})
		return
	}
	request.Description = strings.TrimSpace(request.Description)
	request.URL = strings.TrimSpace(request.URL)
	u, err := url.Parse(request.URL)
	if len(request.URL) > toolboxURLMaxBytes || strings.ContainsAny(request.URL, "\\ \t\r\n") || err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil {
		c.JSON(http.StatusOK, types.Response{Code: codes.BadRequest, Message: "请输入不含账号密码的完整 HTTP 或 HTTPS 链接（最多 2048 字节）", Data: nil})
		return
	}

	response, err := h.service.AdminAddToolboxTool(ctx, request)
	if err != nil {
		c.JSON(http.StatusOK, types.Response{Code: codes.InternalServerError, Message: "添加工具失败", Data: nil})
		return
	}
	c.JSON(http.StatusOK, types.Response{Code: codes.Success, Message: "", Data: response})
}

func (h *toolboxHandler) AdminUpdateToolboxTool(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	ctx := c.Request.Context()
	request := new(types.AdminUpdateToolboxToolRequest)
	if err := c.ShouldBindJSON(request); err != nil {
		h.logger.Warn("parameter binding error", zap.Error(err))
		message := "无效的请求参数"
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
			field := validationErrors[0]
			if field.Field() == "Name" && field.Tag() == "max" {
				message = "工具名称应为 1 至 30 个字符"
			} else if field.Field() == "Description" && field.Tag() == "max" {
				message = "工具简介最多 100 个字符"
			}
		}
		c.JSON(http.StatusOK, types.Response{Code: codes.BadRequest, Message: message, Data: nil})
		return
	}

	if _, err := idutil.ParseID("toolboxID", request.ID); err != nil {
		c.JSON(http.StatusOK, types.Response{Code: codes.BadRequest, Message: "无效的 ID", Data: nil})
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	if request.Name == "" {
		c.JSON(http.StatusOK, types.Response{Code: codes.BadRequest, Message: "工具名称应为 1 至 30 个字符", Data: nil})
		return
	}
	request.Description = strings.TrimSpace(request.Description)
	request.URL = strings.TrimSpace(request.URL)
	u, err := url.Parse(request.URL)
	if len(request.URL) > toolboxURLMaxBytes || strings.ContainsAny(request.URL, "\\ \t\r\n") || err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil {
		c.JSON(http.StatusOK, types.Response{Code: codes.BadRequest, Message: "请输入不含账号密码的完整 HTTP 或 HTTPS 链接（最多 2048 字节）", Data: nil})
		return
	}

	response, err := h.service.AdminUpdateToolboxTool(ctx, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, types.Response{Code: codes.NotFound, Message: "工具不存在", Data: nil})
			return
		}
		c.JSON(http.StatusOK, types.Response{Code: codes.InternalServerError, Message: "更新工具失败", Data: nil})
		return
	}
	c.JSON(http.StatusOK, types.Response{Code: codes.Success, Message: "", Data: response})
}

func (h *toolboxHandler) AdminReorderToolboxTools(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	ctx := c.Request.Context()
	request := new(types.AdminReorderToolboxToolsRequest)
	if err := c.ShouldBindJSON(request); err != nil {
		h.logger.Warn("parameter binding error", zap.Error(err))
		message := "无效的请求参数"
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
			if field := validationErrors[0]; field.Field() == "IDs" && field.Tag() == "unique" {
				message = "排序列表中有重复的工具"
			}
		}
		c.JSON(http.StatusOK, types.Response{Code: codes.BadRequest, Message: message, Data: nil})
		return
	}

	seen := make(map[uint64]struct{}, len(request.IDs))
	for _, rawID := range request.IDs {
		id, err := idutil.ParseID("toolboxID", rawID)
		if err != nil {
			c.JSON(http.StatusOK, types.Response{Code: codes.BadRequest, Message: "无效的 ID", Data: nil})
			return
		}
		if _, exists := seen[id]; exists {
			c.JSON(http.StatusOK, types.Response{Code: codes.BadRequest, Message: "排序列表中有重复的工具", Data: nil})
			return
		}
		seen[id] = struct{}{}
	}

	if err := h.service.AdminReorderToolboxTools(ctx, request); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, types.Response{Code: codes.NotFound, Message: "工具不存在", Data: nil})
			return
		}
		c.JSON(http.StatusOK, types.Response{Code: codes.InternalServerError, Message: "保存工具排序失败", Data: nil})
		return
	}
	c.JSON(http.StatusOK, types.Response{Code: codes.Success, Message: "", Data: nil})
}

func (h *toolboxHandler) AdminDeleteToolboxTool(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	ctx := c.Request.Context()
	request := new(types.AdminDeleteToolboxToolRequest)
	if err := c.ShouldBindJSON(request); err != nil {
		h.logger.Warn("parameter binding error", zap.Error(err))
		c.JSON(http.StatusOK, types.Response{Code: codes.BadRequest, Message: "无效的请求参数", Data: nil})
		return
	}

	if _, err := idutil.ParseID("toolboxID", request.ID); err != nil {
		c.JSON(http.StatusOK, types.Response{Code: codes.BadRequest, Message: "无效的 ID", Data: nil})
		return
	}

	if err := h.service.AdminDeleteToolboxTool(ctx, request); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, types.Response{Code: codes.NotFound, Message: "工具不存在", Data: nil})
			return
		}
		c.JSON(http.StatusOK, types.Response{Code: codes.InternalServerError, Message: "删除工具失败", Data: nil})
		return
	}
	c.JSON(http.StatusOK, types.Response{Code: codes.Success, Message: "", Data: nil})
}
