# TikTok Backend Go (Minimalist Douyin Server)

<p align="center">
  [![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)](https://go.dev/)
  [![Framework](https://img.shields.io/badge/Framework-Gin-008080?style=flat&logo=gin)](https://gin-gonic.com/)
  [![ORM](https://img.shields.io/badge/ORM-GORM-blue?style=flat)](https://gorm.io/)
  [![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
</p>

<p align="center">
  <a href="README.md">🇨🇳 中文</a> &nbsp;|&nbsp; <a href="README_EN.md">🇺🇸 English</a>
</p>

---

<p align="center">
    <strong>Minimalist TikTok Backend in Gin + GORM · Full ByteDance Youth Camp API Implementation</strong>
</p>

## 🌟 Core Features

* 🎬 **Video Feed Stream & Recommendation**: Fully implements the `/douyin/feed/` endpoint, supporting logged-in and guest users to browse videos with time-descending paginated results and `latest_time` cursor pagination.
* 🔐 **User Management & Authentication**: Includes registration, login, and user info retrieval; based on JWT Token global authentication and BCrypt secure password hashing.
* 📤 **Video Upload & Publishing**: Supports uploading MP4 short videos with built-in static asset web serving and work list management.
* ❤️ **Interaction: Likes & Comments**: Idempotent logic for high-concurrent like/unlike operations, liked video list display, comment creation/deletion, and `mm-dd` structured time formatting.
* 👥 **Social Relations & Messaging**: Following/follower lists, auto-detection of mutual (bidirectional) follow friends, client-side scheduled incremental chat polling, and direct message sending.

---

## 🛠️ Technology Stack

* **Programming Language**: Go 1.22+
* **Web Framework**: Gin Framework
* **ORM & Database**: GORM + Pure-Go SQLite (zero-dependency, out-of-the-box) / MySQL 8.0 compatible
* **Authentication**: JWT (`github.com/golang-jwt/jwt/v5`)
* **Security**: BCrypt (`golang.org/x/crypto/bcrypt`)

---

## 📐 Database Schema Design

The project defines 6 core physical data tables with appropriate indexes established for high-frequency queries:

1. **`users` (User Table)**: Stores account password hashes, biographies, avatars, and count caches.
2. **`videos` (Video Table)**: Stores video playback and cover URLs, titles, and publish times; indexed by `idx_created_at` descending.
3. **`favorites` (Like Table)**: Associates users with liked videos; uses `uk_user_video` composite unique index to prevent duplicate likes.
4. **`comments` (Comment Table)**: Stores video comment content and timestamps; indexed by `idx_video_id` to accelerate comment fetching.
5. **`relations` (Follow Table)**: Stores follower-followee relationships; uses `uk_follower_followee` composite unique index.
6. **`messages` (Message Table)**: Stores direct messages; indexed by `idx_from_to` and `idx_created_at` composite indexes to speed up incremental polling.

---

## 🚀 Quick Start

### 1. Start Backend Service (Zero-Config)

The project ships with Pure-Go database drivers and initial demo data — no external MySQL/Redis configuration is required to run:

```bash
# 1. Clone this repository
git clone https://github.com/MeiSiristhebest/tiktok-backend-go.git
cd tiktok-backend-go

# 2. Compile and run
go run ./cmd/server
```

**Expected output**:
```bash
[GIN-debug] Listening and serving HTTP on :8080
[DB] SQLite database auto-migrated: 6 tables created
[DB] Demo data seeded: 3 users, 5 videos, 12 favorites
Server ready at http://127.0.0.1:8080
```

The service will default to running at `http://127.0.0.1:8080`, and will automatically complete SQLite table creation and demo data seeding.

### 2. Verify the Feed Endpoint with curl (works within 5 minutes)

```bash
curl "http://127.0.0.1:8080/douyin/feed/"
```

**Expected response**:
```json
{
  "status_code": 0,
  "status_msg": "success",
  "video_list": [
    {
      "id": 1,
      "author": { "id": 1001, "name": "demo_user_1", "follow_count": 0, "follower_count": 0 },
      "play_url": "http://127.0.0.1:8080/static/videos/demo1.mp4",
      "cover_url": "http://127.0.0.1:8080/static/covers/demo1.jpg",
      "favorite_count": 12,
      "comment_count": 3
    }
  ]
}
```

---

## 📱 Device Integration with Official Minimalist Douyin APK

1. Download and install the official [Minimalist Douyin App APK](https://bytedance.feishu.cn/docx/NMneddpKCoXZJLxHePUcTzGgnmf).
2. Ensure the phone/emulator and the backend computer are on the same LAN (or use 127.0.0.1 for emulator access).
3. On the App home screen (guest state), **double-click the bottom-right "Me" icon** to open 【Advanced Settings】.
4. Enter your computer's LAN IP as the server prefix, e.g.: `http://192.168.1.X:8080`.
5. After saving successfully, you can seamlessly refresh the Feed video stream, register accounts, like/comment, follow users, and send private messages!

---

## 📋 16-Endpoint API Coverage List

| Module | Endpoint Path | Method | Description | Controller |
| :--- | :--- | :--- | :--- | :--- |
| **Basic** | `/douyin/feed/` | GET | Video feed stream | `controller.Feed` |
| **Basic** | `/douyin/user/register/` | POST | User registration | `controller.Register` |
| **Basic** | `/douyin/user/login/` | POST | User login | `controller.Login` |
| **Basic** | `/douyin/user/` | GET | User info | `controller.UserInfo` |
| **Basic** | `/douyin/publish/action/` | POST | Video publish / upload | `controller.PublishAction` |
| **Basic** | `/douyin/publish/list/` | GET | Publish list | `controller.PublishList` |
| **Interaction** | `/douyin/favorite/action/` | POST | Like / unlike action | `controller.FavoriteAction` |
| **Interaction** | `/douyin/favorite/list/` | GET | Liked videos list | `controller.FavoriteList` |
| **Interaction** | `/douyin/comment/action/` | POST | Comment action | `controller.CommentAction` |
| **Interaction** | `/douyin/comment/list/` | GET | Video comment list | `controller.CommentList` |
| **Social** | `/douyin/relation/action/` | POST | Relation action (follow/unfollow) | `controller.RelationAction` |
| **Social** | `/douyin/relation/follow/list/` | GET | Following list | `controller.FollowList` |
| **Social** | `/douyin/relation/follower/list/` | GET | Follower list | `controller.FollowerList` |
| **Social** | `/douyin/relation/friend/list/` | GET | Friend list (mutual follows) | `controller.FriendList` |
| **Social** | `/douyin/message/chat/` | GET | Chat history / messages | `controller.MessageChat` |
| **Social** | `/douyin/message/action/` | POST | Send message action | `controller.MessageAction` |

---

## 🤝 Contributing

Contributions welcome. Quick flow:

```bash
# 1. Fork → Clone → Branch
git checkout -b feat/your-feature

# 2. Compile passes + format check
go build ./...
go vet ./...

# 3. Run unit tests
go test -v ./...

# 4. Commit and open a PR
git commit -m "feat: your change"
git push origin feat/your-feature
```

**Welcome contribution directions**:
- 🔌 Add Redis mode (replace the SQLite in-memory index layer)
- 🧪 Add controller / service layer unit tests
- 🔄 Upgrade chat from polling to WebSocket push
- ⚡ High-concurrency stress testing and performance optimization (stress test data PRs welcome)

---

## 🔒 Security

| Risk Scenario | Mitigation |
|---------|---------|
| **Plaintext Password Storage** | Uses `golang.org/x/crypto/bcrypt` (cost=12) to hash passwords at registration; plaintext is never stored in the DB |
| **JWT Token Forgery** | `golang-jwt/jwt/v5` signature verification; JWT Secret read from environment variables (default dev value is local-only) |
| **Malicious File Upload Masquerading as Video** | Upload endpoint validates `Content-Type` whitelist (only `video/mp4`) plus file magic number secondary confirmation |
| **Unauthorized Like/Comment/Follow** | All interaction endpoints enforce authentication middleware to validate the current Token user_id against the acting subject |
| **SQL Injection** | All queries use GORM parameter binding or `?` placeholders; **string-concatenated SQL is strictly prohibited** |

**Vulnerability disclosure**: Report security issues directly to **`maox_neta@foxmail.com`** — do not file a public issue. We commit to a **first response within 24 hours**.

---

## 📜 License

This project is licensed under the **MIT License**. See the [LICENSE](LICENSE) file for details.
