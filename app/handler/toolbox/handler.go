package toolbox

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"meta-api/app/service/toolbox"
)

type Handler interface {
	AdminGetToolboxList(c *gin.Context)
	AdminAddToolboxTool(c *gin.Context)
	AdminUpdateToolboxTool(c *gin.Context)
	AdminReorderToolboxTools(c *gin.Context)
	AdminDeleteToolboxTool(c *gin.Context)

	UserGetToolboxList(c *gin.Context)
}

type toolboxHandler struct {
	logger  *zap.Logger
	service toolbox.Service
}

func NewHandler(logger *zap.Logger, service toolbox.Service) Handler {
	return &toolboxHandler{
		logger:  logger,
		service: service,
	}
}
