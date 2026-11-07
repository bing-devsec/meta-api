---
name: blog-api
description: 修改 meta-api 的 Go 接口、鉴权、文章发布、评论审核或 Redis 缓存时使用，保持管理端与前台的现有协议和数据一致性。仅修改 Nuxt 页面或本地 JSON 编辑器时不使用。
---

# 博客后端开发

本技能的源码路径均相对于 `meta-api/` 仓库根目录。

## 先定位实际链路

- 主后端是本仓库；同级 `api-server` 是旧工程，不把修改顺带同步到那里。
- 从 `app/router/user_routes.go` 或 `admin_routes.go` 找入口，再追踪 handler、service、model。
- Handler 负责绑定与响应，service 编排业务，model 负责 MySQL 读写；新增依赖由 `app/di` 统一装配。
- 请求/响应结构放在 `common/types`，业务码来自 `common/codes`，Redis key 由 `common/cachekey` 构造。
- 框架、依赖与工具链以当前 `go.mod` 为准，先检查工作区已有改动；设计文档不覆盖源码事实。

## 保持接口契约

- 现有业务接口通常返回 HTTP 200 和 `{code, message, data}`；成功码为 `2000`，成功消息通常为空。
- 沿用相邻接口的错误映射，不能把 HTTP 成功当成业务成功，也不批量改写现有状态码语义。
- Sonyflake ID 在业务内为 `uint64`，在 HTTP 协议中使用字符串；解析使用 `common/idutil.ParseID`。
- 管理端 `/api` 反向代理前缀不属于 Gin 路由；后端注册 `/admin` 和 `/user`。
- 管理员登录、普通用户 OAuth、匿名 guard 是不同边界；guard 不能替代登录与权限校验。
- 修改认证、Cookie、超时或代理信任时，阅读 [认证与验证](references/auth-and-verification.md)。

## 维护写入与缓存边界

- 涉及文章/草稿、图片引用、评论可见性或排序缓存时，阅读 [内容生命周期](references/content-lifecycle.md)。
- 数据库事务只保证 MySQL 内的原子性；Redis、CDN、COS、sitemap 调用有独立失败窗口。
- 修改多阶段写入时明确已提交的数据、外部失败如何返回、重试是否会重复操作，不把它们假定为一笔事务。
- 保留请求上下文和取消语义；新增接口不顺便放大全局超时或可信代理范围。
- 使用 `WithContext(ctx)` 的现有模型方法，保持单数表名与已发布/草稿筛选条件。
- 修改模型前检查 `bootstrap/mysql.go` 的 AutoMigrate 调用和实际使用条件，说明 schema 影响。

## 验证与交付

- 根据变更选择受影响单元测试、Go 构建/静态检查，以及必要的本地接口回归；不为文档变更启动服务。
- 执行命令前确认当前测试文件和测试名称确实存在，不能把 `no tests to run` 当成覆盖通过。
- 后端测试现状、纯编译命令、需服务联调的边界见 [认证与验证](references/auth-and-verification.md)。
- 不用 `go run .`、重启现有后端或调用写接口来证明代码能编译；应用启停具有真实持久化副作用。
- 联调写操作使用隔离测试数据与已确认的本地环境，覆盖实际改变的成功、失败和权限分支。
- 交付说明接口行为、兼容影响、实际验证结果，以及未验证的外部依赖；不把编译通过描述成端到端通过。

## 授权与配置

- 常规源码修改、构建与隔离本地验证按当前任务授权推进，无需额外增加批准步骤。
- 部署、生产重启、生产 schema/内容变更、CDN purge 与 COS 删除仅在当前用户授权涵盖时执行。
- 加载本技能不扩大授权；核对运行环境、操作对象和现有范围后再执行外部变更。
- 非敏感业务配置、运行时环境与 secrets 的来源分层见 `main.go`、`bootstrap/env.go`。
- 调试报告只保留必要的变量名与脱敏错误，不把密钥、Cookie、数据库连接串写进日志或提交。

## 按需查阅

- 架构和生命周期：`docs/后端代码架构.md`、`docs/数据库设计.md`。
- Redis 与浏览量：`docs/缓存设计.md`、`docs/文章浏览量设计.md`。
- 审核与分享：`docs/评论审核系统设计.md`、`docs/结构化数据分享设计.md`。
- 这些文档可能含过时路径或测试命令，使用前以当前文件、调用点和工作区状态核实。
