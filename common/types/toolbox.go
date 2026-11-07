package types

type AdminToolboxItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	SortOrder   int    `json:"sortOrder"`
	Enabled     bool   `json:"enabled"`
}

type UserToolboxItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	SortOrder   int    `json:"sortOrder"`
	Enabled     bool   `json:"enabled"`
}

type AdminGetToolboxListResponse struct {
	Tools []AdminToolboxItem `json:"tools"`
}

type UserGetToolboxListResponse struct {
	Tools []UserToolboxItem `json:"tools"`
}

type AdminAddToolboxToolRequest struct {
	Name        string `json:"name" binding:"required,max=30"`
	Description string `json:"description" binding:"max=100"`
	URL         string `json:"url" binding:"required"`
	// 省略时追加到末尾；显式值兼容旧客户端。
	SortOrder *int  `json:"sortOrder" binding:"omitempty,gte=0,lte=1000000"`
	Enabled   *bool `json:"enabled" binding:"required"`
}

type AdminAddToolboxToolResponse struct {
	ID string `json:"id"`
}

type AdminUpdateToolboxToolRequest struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required,max=30"`
	Description string `json:"description" binding:"max=100"`
	URL         string `json:"url" binding:"required"`
	// 省略时保留现有顺序。
	SortOrder *int  `json:"sortOrder" binding:"omitempty,gte=0,lte=1000000"`
	Enabled   *bool `json:"enabled" binding:"required"`
}

type AdminUpdateToolboxToolResponse struct {
	ID string `json:"id"`
}

type AdminReorderToolboxToolsRequest struct {
	IDs []string `json:"ids" binding:"required,unique,dive,required"`
}

type AdminDeleteToolboxToolRequest struct {
	ID string `json:"id" binding:"required"`
}
