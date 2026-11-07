package toolbox

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"

	toolboxModel "meta-api/app/model/toolbox"
	"meta-api/common/types"
)

type Service interface {
	AdminGetToolboxList(ctx context.Context) (*types.AdminGetToolboxListResponse, error)
	AdminAddToolboxTool(ctx context.Context, request *types.AdminAddToolboxToolRequest) (*types.AdminAddToolboxToolResponse, error)
	AdminUpdateToolboxTool(ctx context.Context, request *types.AdminUpdateToolboxToolRequest) (*types.AdminUpdateToolboxToolResponse, error)
	AdminReorderToolboxTools(ctx context.Context, request *types.AdminReorderToolboxToolsRequest) error
	AdminDeleteToolboxTool(ctx context.Context, request *types.AdminDeleteToolboxToolRequest) error

	UserGetToolboxList(ctx context.Context) (*types.UserGetToolboxListResponse, error)
}

type toolboxService struct {
	logger      *zap.Logger
	idGenerator *sonyflake.Sonyflake
	redis       *redis.Client
	model       toolboxModel.Model
}

func NewService(logger *zap.Logger, idGenerator *sonyflake.Sonyflake, redis *redis.Client, model toolboxModel.Model) Service {
	return &toolboxService{
		logger:      logger,
		idGenerator: idGenerator,
		redis:       redis,
		model:       model,
	}
}
