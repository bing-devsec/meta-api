# 认证、运行环境与验证

仅在认证、网络边界、后端验证或服务联调相关任务中阅读。源码路径相对于 `meta-api/` 根目录。

## 身份边界

- 管理员路由组：`app/router/admin_routes.go`，保护入口为 `/admin/auth/*`。
- 管理员中间件：`common/middlewares/admin_login_jwt.go`；令牌从 Cookie 读取。
- 普通用户评论身份：`common/middlewares/comment_user_jwt.go`、`app/service/userauth`。
- OAuth 的安全单元测试：`app/service/userauth/oauth_security_test.go`。
- 匿名 guard：`common/guard`；它不提供管理员或评论用户权限。

保持管理员与普通用户 token、Cookie、session version 的不同用途。前端 `localStorage.userID` 只是登录界面标记，服务端始终负责权限判断。

管理员 Cookie 工具在 `common/utils/cookie.go`：HttpOnly；生产 Secure/Strict，本地 HTTP 使用 Lax。refresh Cookie 的 Path 包含浏览器 `/api` 前缀，删除时 Path 必须一致。不要为解决本地代理问题而削弱生产 Cookie 设置。

OAuth 保留站内回跳校验、浏览器 flow binding、state 消费顺序和 PKCE；修改后运行实际存在的安全测试。回调异常不能记录令牌、Cookie 或完整密钥材料。

## HTTP 边界

`app/router/router_engine.go` 目前默认请求超时 3 秒，对 OAuth、文章写操作等有路由级覆盖。调整慢接口时优先核对请求路径和相应覆盖，不直接延长全部请求。

Gin 只信任指定 Docker 代理网段传来的客户端地址头；不要引入绕过网段校验的全局头信任。管理端本地代理见同级 `admin-ui/vite.config.ts`，线上 Nginx/CDN 规则不能仅凭本仓库文档认定已配置。

## 无服务验证

先检查 `go.mod`、`git status --short`、现存 `*_test.go` 和相关包的初始化副作用，再选命令。常用入口：

```sh
go test ./app/service/userauth ./app/model/comment
go test ./...
go vet ./...
go build ./...
```

按改动风险选择，不要求每次都全部运行。现有单元测试覆盖 OAuth 工具与审核表名，不代表文章、Redis、CDN 或真实登录均已回归。`go test` 提示没有匹配测试时应写明未覆盖。

审核文档中引用的 `bootstrap/config_test.go`、`bootstrap/moderation_corpus_test.go` 可能已不在当前工作区；先查文件和 `go test -list` 的实际输出，不盲用文档中的 `-run` 命令，也不要擅自恢复用户删除的测试。

新增测试优先针对可观察协议、安全边界或真实回归路径；不要只断言实现细节或为了复述技能编写测试。模型或服务测试可注入假依赖，避免借用真实数据库、Redis 或云服务。

## 需要服务的验证

`main.go` 的启动会加载本地环境或生产 secrets，连接 MySQL/Redis；`bootstrap/app.go` 的 `InitMySQL` 当前会调用 `bootstrap/mysql.go` 中的 AutoMigrate。`app/app.go` 会预热缓存、启动 cron；停止阶段会将浏览量写回 MySQL。

因此不以 `go run .` 或重启现有服务替代编译验证。确需联调时先确认目标确为隔离本地实例、数据可丢弃、外部客户端不会指向生产，再按当前用户授权启动或复用服务；不要无故接管已有进程。

测试接口同时判断 HTTP 状态和响应业务码；授权失败、token 过期与真实服务故障是不同结果。写入联调按任务覆盖成功、输入错误、无权限、重复提交及外部失败中相关分支，不机械要求全部场景。

## 配置与发布范围

业务 YAML、非敏感环境变量和 secrets 分层由 `main.go`、`bootstrap/env.go`、`common/utils` 定义。审查环境变量名和来源即可定位多数配置问题，不把 `.env` 内容、Cookie 或连接串复制到产物。

本地构建无需新的批准流程。部署、生产重启、schema 变更、CDN purge 和 COS 删除是独立的外部动作：仅当用户当前授权已覆盖时执行，并核对目标；授权缺失时完成可审阅的代码和验证后再请求所缺授权。
