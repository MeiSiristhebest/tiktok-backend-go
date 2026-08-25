# ⚡ tiktok-backend-go Architecture Blueprint

<p align="center">
  <b>English | <a href="./ARCHITECTURE_zh.md">简体中文</a></b>
</p>

This document details the high-concurrency video feed, user interaction, and relation graph architecture powering **tiktok-backend-go**.

```mermaid
graph TD
    Client[Mobile App / Web Client] -->|HTTP REST / JWT| Router[Gin Engine Router]

    subgraph "Middleware Pipeline"
        Router --> Cors[CORS & Recovery]
        Router --> JWTAuth[JWT Auth Middleware]
        Router --> RateLimit[Token Bucket Limiter]
    end

    subgraph "Core Business Handlers"
        JWTAuth --> FeedHandler[Feed & Video Streaming]
        JWTAuth --> UserHandler[User & Profile Service]
        JWTAuth --> RelationHandler[Follow / Follower Graph]
        JWTAuth --> CommentHandler[Comment & Favorite Service]
    end

    subgraph "Storage & Caching"
        FeedHandler --> Redis[(Redis Cache Cluster)]
        RelationHandler --> Redis
        FeedHandler --> MySQL[(MySQL InnoDB Storage)]
        UserHandler --> MySQL
        RelationHandler --> MySQL
    end
```

---

## 🚀 1. Video Feed Caching & Low-Latency Delivery
- Utilizes Redis sorted sets (`ZSET`) keyed by publish timestamp to power sub-millisecond timeline feed pagination.
- Asynchronously increments video view counts and interaction counters with batched Redis flushes to MySQL.

---

## 🔒 2. Relation Graph Transaction Handling
- Atomically updates followings and follower counts within ACID database transactions.
- Employs Redis bitmap & set operations to achieve O(1) mutual-follow relation queries.

---

<sub>© 2026 tiktok-backend-go. Licensed under the MIT License.</sub>
