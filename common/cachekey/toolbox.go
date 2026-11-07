package cachekey

// ToolboxEnabledList 前台可见工具的完整列表，按后台拖拽顺序排列。
func ToolboxEnabledList() Key { return build("toolbox", "enabled", "list") }
