package toolbox

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Toolbox struct {
	ID          uint64    `gorm:"column:id;type:bigint unsigned;primary_key;autoIncrement:false;NOT NULL;index:idx_toolbox_sort,priority:2;index:idx_toolbox_enabled_sort,priority:3"`
	Name        string    `gorm:"column:name;type:varchar(80);NOT NULL"`
	Description string    `gorm:"column:description;type:varchar(300);NOT NULL"`
	URL         string    `gorm:"column:url;type:varchar(2048);NOT NULL"`
	SortOrder   int       `gorm:"column:sort_order;type:bigint;NOT NULL;index:idx_toolbox_sort,priority:1;index:idx_toolbox_enabled_sort,priority:2"`
	Enabled     bool      `gorm:"column:enabled;type:boolean;NOT NULL;index:idx_toolbox_enabled_sort,priority:1"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime(3);NOT NULL;default:CURRENT_TIMESTAMP(3)"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime(3);NOT NULL;default:CURRENT_TIMESTAMP(3)"`
}

func (m *toolboxModel) ListToolboxes(ctx context.Context) ([]Toolbox, error) {
	rows := make([]Toolbox, 0)
	if err := m.mysql.WithContext(ctx).Model(&Toolbox{}).
		Order("sort_order ASC").
		Order("id ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list toolbox entries: %w", err)
	}
	return rows, nil
}

func (m *toolboxModel) ListEnabledToolboxes(ctx context.Context) ([]Toolbox, error) {
	rows := make([]Toolbox, 0)
	if err := m.mysql.WithContext(ctx).Model(&Toolbox{}).
		Where("enabled = ?", true).
		Order("sort_order ASC").
		Order("id ASC").
		Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list enabled toolbox entries: %w", err)
	}
	return rows, nil
}

// CreateToolbox 在未指定排序值时追加到列表末尾。
func (m *toolboxModel) CreateToolbox(ctx context.Context, item *Toolbox, sortOrder *int) error {
	if sortOrder != nil {
		item.SortOrder = *sortOrder
	} else {
		var maxSortOrder int
		if err := m.mysql.WithContext(ctx).Model(&Toolbox{}).
			Select("COALESCE(MAX(sort_order), -1)").Scan(&maxSortOrder).Error; err != nil {
			return fmt.Errorf("get max toolbox sort order: %w", err)
		}
		item.SortOrder = maxSortOrder + 1
	}
	if err := m.mysql.WithContext(ctx).Model(&Toolbox{}).Create(item).Error; err != nil {
		return fmt.Errorf("create toolbox: %w", err)
	}
	return nil
}

// ReorderToolboxes 按传入的 ID 顺序，一次更新所有排序值。
func (m *toolboxModel) ReorderToolboxes(ctx context.Context, ids []uint64) error {
	if len(ids) == 0 {
		return nil
	}
	var count int64
	if err := m.mysql.WithContext(ctx).Model(&Toolbox{}).Where("id IN ?", ids).Count(&count).Error; err != nil {
		return fmt.Errorf("check toolbox IDs for reorder: %w", err)
	}
	if count != int64(len(ids)) {
		return gorm.ErrRecordNotFound
	}

	cases := make([]string, len(ids))
	args := make([]any, 0, len(ids)*2)
	for index, id := range ids {
		cases[index] = "WHEN ? THEN ?"
		args = append(args, id, index)
	}
	sortOrder := gorm.Expr("CASE id "+strings.Join(cases, " ")+" END", args...)
	if err := m.mysql.WithContext(ctx).Model(&Toolbox{}).
		Where("id IN ?", ids).UpdateColumn("sort_order", sortOrder).Error; err != nil {
		return fmt.Errorf("reorder toolboxes: %w", err)
	}
	return nil
}

func (m *toolboxModel) UpdateToolbox(ctx context.Context, item *Toolbox, sortOrder *int) error {
	fields := []string{"name", "description", "url", "enabled", "update_time"}
	if sortOrder != nil {
		item.SortOrder = *sortOrder
		fields = append(fields, "sort_order")
	}
	result := m.mysql.WithContext(ctx).Model(&Toolbox{}).Where("id = ?", item.ID).Select(fields).Updates(item)
	if result.Error != nil {
		return fmt.Errorf("update toolbox: %w", result.Error)
	}
	// MySQL 对内容未变化的更新也可能返回 0，不能直接判定记录不存在。
	if result.RowsAffected == 0 {
		var count int64
		if err := m.mysql.WithContext(ctx).Model(&Toolbox{}).Where("id = ?", item.ID).Count(&count).Error; err != nil {
			return fmt.Errorf("confirm toolbox exists: %w", err)
		}
		if count == 0 {
			return gorm.ErrRecordNotFound
		}
	}
	return nil
}

func (m *toolboxModel) DeleteToolbox(ctx context.Context, id uint64) error {
	result := m.mysql.WithContext(ctx).Model(&Toolbox{}).Where("id = ?", id).Delete(&Toolbox{})
	if result.Error != nil {
		return fmt.Errorf("delete toolbox: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
