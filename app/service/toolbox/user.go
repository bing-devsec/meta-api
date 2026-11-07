package toolbox

import (
	"context"
	"errors"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"meta-api/common/cachekey"
	"meta-api/common/types"
)

func (s *toolboxService) UserGetToolboxList(ctx context.Context) (*types.UserGetToolboxListResponse, error) {
	key := cachekey.ToolboxEnabledList().String()
	if s.redis != nil {
		value, err := s.redis.Get(ctx, key).Bytes()
		if err == nil {
			items := make([]types.UserToolboxItem, 0)
			if err = sonic.Unmarshal(value, &items); err == nil && items != nil {
				return &types.UserGetToolboxListResponse{Tools: items}, nil
			}
		}
		if err != nil && !errors.Is(err, redis.Nil) {
			s.logger.Warn("failed to read toolbox cache", zap.Error(err))
		}
	}

	rows, err := s.model.ListEnabledToolboxes(ctx)
	if err != nil {
		s.logger.Error("failed to list enabled toolbox entries", zap.Error(err))
		return nil, err
	}
	items := make([]types.UserToolboxItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, types.UserToolboxItem{
			ID:          strconv.FormatUint(row.ID, 10),
			Name:        row.Name,
			Description: row.Description,
			URL:         row.URL,
			SortOrder:   row.SortOrder,
			Enabled:     row.Enabled,
		})
	}
	if s.redis != nil {
		value, err := sonic.Marshal(items)
		if err != nil {
			return nil, err
		}
		// 包括空列表在内，公开业务缓存均不设置过期时间。
		if err = s.redis.Set(ctx, key, value, 0).Err(); err != nil {
			s.logger.Warn("failed to cache enabled toolbox list", zap.Error(err))
		}
	}
	return &types.UserGetToolboxListResponse{Tools: items}, nil
}
