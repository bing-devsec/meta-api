# 内容生命周期与一致性

仅在变更文章、草稿、图片、评论或相关缓存时阅读。源码路径相对于 `meta-api/` 根目录。

## 文章与草稿

- 路由：`app/router/admin_routes.go`。
- HTTP 输入输出：`app/handler/article`、`common/types/article.go`。
- 编排：`app/service/article/admin.go`、`draft.go`、`image.go`。
- 持久化：`app/model/article/article.go`、`draft.go`、`image.go`。

文章和草稿共用文章表，以 `status` 和 `published_id` 区分；公开查询必须限定已发布状态。保存编辑草稿不能改变当前公开文章。发布新文章与发布已有文章的编辑草稿是两条不同分支，不在重构中合并掉它们的 ID、时间和浏览量语义。

已有文章更新需要保留 Redis 中最新的浏览量，不能用旧 MySQL 读数覆盖。发布还会同步图片引用、文章/标签缓存和 sitemap；更新已有文章需要清理文章 HTML 的 CDN 缓存。具体顺序按当前实现核实，不将多系统操作描述成全局原子事务。

当前删除文章先调用 CDN purge；失败时保留数据库文章。修改此顺序会改变错误后的业务状态，需要给出明确理由与回归证据。模型中的文章/关联草稿事务不自动涵盖 Redis 或外部客户端。

遇到“后台提示失败但已经提交”的路径，先定位失败发生在数据库提交前还是提交后，再决定返回与重试行为。不得盲目重新执行新增、发布、图片上传或删除。

## Redis 与浏览量

- 排序和预热：`app/service/article/cache.go`、`list_cache.go`。
- 标签缓存：`app/service/tag/article_cache.go`。
- Key 构造：`common/cachekey/article.go`、`tag.go`、`comment.go`。
- 打点计数：`app/service/viewlog`。

文章排序使用时间/浏览量 ZSet。预热先分批构造临时 key，再通过事务切换正式 key；切换后需去除临时 TTL，保留空集合哨兵的处理。不要改成先删除正式 key 再回填，避免空窗口或半成品可见。

文章更新会影响详情 Hash、标签计数和新旧标签列表。变更排序或标签行为时检查缓存命中、缺失回源、空结果和写入后的失效；不要只测试一次冷读取。使用现有锁和 key 工厂，避免另造不同命名的缓存副本。

浏览量以 Redis 为实时来源，定时任务和停机阶段回写 MySQL。重启服务可能引起实际写入；修改此链路时核对 HTTP 停止、cron 停止与最终持久化顺序。

## 评论与审核

- 提交编排：`app/service/comment/user.go`。
- 审核：`app/service/comment/moderation`、`moderation_audit.go`。
- 缓存：`app/service/comment/cache.go`。
- 审计持久化：`app/model/comment/moderation_audit.go`。
- 策略入口：`config/comment_moderation.manifest.yml`，细分配置在 `config/moderation/`。

仅 `approved` 评论可公开；`pending`、`rejected` 不能混入公开列表。审核错误保持待审核的安全语义。普通用户响应不能顺便暴露内部拒绝原因、策略证据或审计轨迹。

真实评论及审核快照同事务写入；持久化成功后才记录行为统计。评论通过、状态改变、举报下架、删除，以及影响公开展示的用户资料变化，都需要核对相应文章评论快照的失效。

修改策略时同时检验应拦截样例、正常表达和灰区，不能只添加更多拒绝规则。`docs/评论审核系统设计.md` 提到的黄金语料测试先确认仍存在；缺失测试入口时明确记录覆盖缺口，不将空测试当成策略回归。

## 管理端与前台协同

响应字段变化需要核对同级 `admin-ui` 与 `portal-web` 消费方。内容发布回归使用隔离数据验证草稿保存、发布后的公开读取及必要缓存更新，不能用生产文章作可丢弃样例。

图片资产可能调用 COS，更新/删除可能调用 CDN，sitemap 刷新也有外部边界。无测试凭据时使用 fake/桩或报告未验证部分，不擅自改用生产凭据。
