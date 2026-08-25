# ⚡ tiktok-backend-go 高并发短视频服务架构设计 (Architecture Guide)

<p align="center">
  <b><a href="./ARCHITECTURE.md">English</a> | 简体中文</b>
</p>

本文档阐述 **tiktok-backend-go** 在高并发短视频流推送、用户交互与社交关系链场景下的系统架构与优化实践。

```mermaid
graph TD
    Client[移动端 App / Web 客户端] -->|HTTP REST / JWT| Router[Gin 高性能路由引擎]

    subgraph "中间件流水线"
        Router --> Cors[跨域与异常捕获 Recovery]
        Router --> JWTAuth[JWT 统一鉴权中间件]
        Router --> RateLimit[令牌桶流量控制]
    end

    subgraph "核心业务层"
        JWTAuth --> FeedHandler[Feed 视频流处理]
        JWTAuth --> UserHandler[用户中心服务]
        JWTAuth --> RelationHandler[关注与粉丝关系链]
        JWTAuth --> CommentHandler[点赞与评论服务]
    end

    subgraph "存储与多级缓存"
        FeedHandler --> Redis[(Redis 高速缓存集群)]
        RelationHandler --> Redis
        FeedHandler --> MySQL[(MySQL 8.0 关系数据库)]
        UserHandler --> MySQL
        RelationHandler --> MySQL
    end
```

---

## 🚀 1. 毫秒级短视频 Feed 流缓存架构
- 借助 Redis `ZSET` 基于发布时间戳维护视频索引流，实现毫秒级分页与滑动窗口查询。
- 点赞与播放计数器在 Redis 中进行高并发自增，并通过后台定时任务批量刷盘至 MySQL。

---

## 🔒 2. 社交关系链事务一致性
- 关注/取关操作采用 MySQL 事务保障粉丝数与关注数双向一致。
- 借助 Redis `SET` 的交集计算实现 O(1) 复杂度的互相关注/好友关系极速判定。

---

<sub>© 2026 tiktok-backend-go. Licensed under the MIT License.</sub>
