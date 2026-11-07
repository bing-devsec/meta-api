package toolbox

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	toolboxModel "meta-api/app/model/toolbox"
	"meta-api/common/cachekey"
	"meta-api/common/idutil"
	"meta-api/common/types"
)

func (s *toolboxService) AdminGetToolboxList(ctx context.Context) (*types.AdminGetToolboxListResponse, error) {
	rows, err := s.model.ListToolboxes(ctx)
	if err != nil {
		s.logger.Error("failed to list toolbox entries", zap.Error(err))
		return nil, err
	}
	items := make([]types.AdminToolboxItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, types.AdminToolboxItem{
			ID:          strconv.FormatUint(row.ID, 10),
			Name:        row.Name,
			Description: row.Description,
			URL:         row.URL,
			SortOrder:   row.SortOrder,
			Enabled:     row.Enabled,
		})
	}
	return &types.AdminGetToolboxListResponse{Tools: items}, nil
}

func (s *toolboxService) AdminAddToolboxTool(ctx context.Context, request *types.AdminAddToolboxToolRequest) (*types.AdminAddToolboxToolResponse, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		s.logger.Error("failed to get toolbox entry time", zap.Error(err))
		return nil, fmt.Errorf("load location: %w", err)
	}
	now := time.Now().In(loc)
	id, err := s.idGenerator.NextID()
	if err != nil {
		s.logger.Error("failed to generate toolbox entry id", zap.Error(err))
		return nil, fmt.Errorf("generate toolbox entry id: %w", err)
	}
	item := &toolboxModel.Toolbox{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
		URL:         request.URL,
		Enabled:     *request.Enabled,
		CreateTime:  now,
		UpdateTime:  now,
	}
	if err = s.model.CreateToolbox(ctx, item, request.SortOrder); err != nil {
		s.logger.Error("failed to create toolbox entry", zap.Error(err))
		return nil, err
	}
	if err := s.clearEnabledCache(ctx); err != nil {
		return nil, err
	}
	return &types.AdminAddToolboxToolResponse{ID: strconv.FormatUint(id, 10)}, nil
}

func (s *toolboxService) AdminUpdateToolboxTool(ctx context.Context, request *types.AdminUpdateToolboxToolRequest) (*types.AdminUpdateToolboxToolResponse, error) {
	id, err := idutil.ParseID("toolboxID", request.ID)
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		s.logger.Error("failed to get toolbox entry time", zap.Error(err))
		return nil, fmt.Errorf("load location: %w", err)
	}
	item := &toolboxModel.Toolbox{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
		URL:         request.URL,
		Enabled:     *request.Enabled,
		UpdateTime:  time.Now().In(loc),
	}
	if err = s.model.UpdateToolbox(ctx, item, request.SortOrder); err != nil {
		s.logger.Error("failed to update toolbox entry", zap.Error(err))
		return nil, err
	}
	if err := s.clearEnabledCache(ctx); err != nil {
		return nil, err
	}
	return &types.AdminUpdateToolboxToolResponse{ID: strconv.FormatUint(id, 10)}, nil
}

func (s *toolboxService) AdminReorderToolboxTools(ctx context.Context, request *types.AdminReorderToolboxToolsRequest) error {
	ids := make([]uint64, 0, len(request.IDs))
	for _, rawID := range request.IDs {
		id, err := idutil.ParseID("toolboxID", rawID)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}
	if err := s.model.ReorderToolboxes(ctx, ids); err != nil {
		s.logger.Error("failed to reorder toolbox entries", zap.Error(err))
		return err
	}
	if err := s.clearEnabledCache(ctx); err != nil {
		return err
	}
	return nil
}

func (s *toolboxService) AdminDeleteToolboxTool(ctx context.Context, request *types.AdminDeleteToolboxToolRequest) error {
	id, err := idutil.ParseID("toolboxID", request.ID)
	if err != nil {
		return err
	}
	if err = s.model.DeleteToolbox(ctx, id); err != nil {
		s.logger.Error("failed to delete toolbox entry", zap.Error(err))
		return err
	}
	if err := s.clearEnabledCache(ctx); err != nil {
		return err
	}
	return nil
}

func (s *toolboxService) clearEnabledCache(ctx context.Context) error {
	if s.redis == nil {
		return nil
	}
	if err := s.redis.Del(ctx, cachekey.ToolboxEnabledList().String()).Err(); err != nil {
		s.logger.Error("MySQL mutation succeeded but failed to clear toolbox cache", zap.Error(err))
		return fmt.Errorf("MySQL mutation succeeded but failed to clear toolbox cache: %w", err)
	}
	return nil
}
