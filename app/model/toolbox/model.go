package toolbox

import (
	"context"

	"gorm.io/gorm"
)

type Model interface {
	ListToolboxes(ctx context.Context) ([]Toolbox, error)
	ListEnabledToolboxes(ctx context.Context) ([]Toolbox, error)
	CreateToolbox(ctx context.Context, item *Toolbox, sortOrder *int) error
	UpdateToolbox(ctx context.Context, item *Toolbox, sortOrder *int) error
	DeleteToolbox(ctx context.Context, id uint64) error
	ReorderToolboxes(ctx context.Context, ids []uint64) error
}

type toolboxModel struct {
	mysql *gorm.DB
}

func NewModel(mysql *gorm.DB) Model {
	return &toolboxModel{mysql: mysql}
}
