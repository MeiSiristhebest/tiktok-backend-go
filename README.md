# TikTok Backend Go (极简版抖音服务端)

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![Framework](https://img.shields.io/badge/Framework-Gin-008080?style=flat&logo=gin)](https://gin-gonic.com)
[![ORM](https://img.shields.io/badge/ORM-GORM-blue?style=flat)](https://gorm.io)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

> **Project Overview**  
> 本项目基于**字节跳动第六届青训营**后端项目需求规范进行架构设计与工程落地方案开发，使用 **Golang (Gin + GORM + MySQL / SQLite)** 打造极简版抖音 Backend 服务。
> 
> 全量实现了青训营要求的 **基础服务、互动服务、社交服务** 3 大方向共 16 个 RESTful API 接口，并配合官方极简抖音 App (APK) 进行了全流程功能联调与性能验证。

---

## 🌟 核心特性 (Features)

* 🎬 **视频 Feed 流与推荐**：全量实现 `/douyin/feed/` 接口，支持未登录/登录用户刷视频，按投稿时间倒序推出的限制分页与 `latest_time` 翻页。
* 🔐 **用户与鉴权**：包含注册、登录与用户信息拉取，基于 JWT Token 全局鉴权与 BCrypt 密码安全哈希存储。
* 📤 **视频投稿与播放**：支持上传 MP4 短视频，内置静态资源 Web 服务与作品列表管理。
* ❤️ **互动点赞与评论**：高并发点赞/取消点赞幂等逻辑、喜欢（点赞）视频列表展示、评论发表/删除及 `mm-dd` 结构化时间显示。
* 👥 **社交关系链与聊天**：关注/粉丝列表、好友（双向关注）自动化判别，客户端定时轮询增量聊天记录与私信发送。

---

## 🛠️ 技术选型 (Tech Stack)

* **Programming Language**: Go 1.22+
* **Web Framework**: Gin Framework
* **ORM & Database**: GORM + Pure-Go SQLite (零依赖开箱即用) / MySQL 8.0 兼容
* **Authentication**: JWT (`github.com/golang-jwt/jwt/v5`)
* **Security**: BCrypt (`golang.org/x/crypto/bcrypt`)

---

## 📐 数据库架构设计 (Database Schema)

项目设计了 6 张物理核心数据表，并针对高频查询建立了合理的索引：

1. **`users` (用户表)**：保存账号密码哈希、个性签名、头像及统计基数。
2. **`videos` (视频表)**：保存视频播放与封面 URL、标题及发布时间，建有 `idx_created_at` 倒序索引。
3. **`favorites` (点赞表)**：用户与视频的点赞关联，建立 `uk_user_video` 联合唯一索引防止重复点赞。
4. **`comments` (评论表)**：视频评论内容及时间，建有 `idx_video_id` 提高评论拉取速度。
5. **`relations` (关注表)**：粉丝与被关注者关系，建立 `uk_follower_followee` 联合唯一索引。
6. **`messages` (消息表)**：私信记录，建有 `idx_from_to` 及 `idx_created_at` 复合索引加速增量轮询。

---

## 🚀 快速开始 (Quick Start)

### 1. 启动后端服务 (Zero-Config)

项目内置 Pure-Go 数据库驱动与初始测试数据，无需额外配置 MySQL/Redis 即可直接运行：

```bash
# 1. 克隆本仓库
git clone https://github.com/MeiSiristhebest/tiktok-backend-go.git
cd tiktok-backend-go

# 2. 编译并运行
go run ./cmd/server
```

服务启动后将默认运行在 `http://127.0.0.1:8080`，并自动完成 SQLite 数据库建表与演示数据填充。

---

## 📱 与极简抖音 APK 客户端实机联调

1. 下载并安装官方 [极简抖音 App APK](https://bytedance.feishu.cn/docx/NMneddpKCoXZJLxHePUcTzGgnmf)。
2. 保证手机/模拟器与运行后端的电脑处于同一局域网（或使用 127.0.0.1 模拟器访问）。
3. 打开 App 首页，未登录状态下**双击右下角 “我”** 打开【高级设置】。
4. 在服务器前缀中填入你的电脑 IP，例如：`http://192.168.1.X:8080`。
5. 成功保存后即可顺畅刷新 Feed 视频流、注册账号、点赞评论、关注与发送私信！

---

## 📋 16 个 API 接口覆盖清单

| 模块 | 接口路径 | 类型 | 说明 | 对应 Controller |
| :--- | :--- | :--- | :--- | :--- |
| **基础** | `/douyin/feed/` | GET | 视频流接口 | `controller.Feed` |
| **基础** | `/douyin/user/register/` | POST | 用户注册 | `controller.Register` |
| **基础** | `/douyin/user/login/` | POST | 用户登录 | `controller.Login` |
| **基础** | `/douyin/user/` | GET | 用户信息 | `controller.UserInfo` |
| **基础** | `/douyin/publish/action/` | POST | 视频投稿 | `controller.PublishAction` |
| **基础** | `/douyin/publish/list/` | GET | 发布列表 | `controller.PublishList` |
| **互动** | `/douyin/favorite/action/` | POST | 赞操作 | `controller.FavoriteAction` |
| **互动** | `/douyin/favorite/list/` | GET | 喜欢列表 | `controller.FavoriteList` |
| **互动** | `/douyin/comment/action/` | POST | 评论操作 | `controller.CommentAction` |
| **互动** | `/douyin/comment/list/` | GET | 视频评论列表 | `controller.CommentList` |
| **社交** | `/douyin/relation/action/` | POST | 关系操作 (关注/取关) | `controller.RelationAction` |
| **社交** | `/douyin/relation/follow/list/` | GET | 关注列表 | `controller.FollowList` |
| **社交** | `/douyin/relation/follower/list/` | GET | 粉丝列表 | `controller.FollowerList` |
| **社交** | `/douyin/relation/friend/list/` | GET | 好友列表 | `controller.FriendList` |
| **社交** | `/douyin/message/chat/` | GET | 聊天记录 | `controller.MessageChat` |
| **社交** | `/douyin/message/action/` | POST | 消息操作 | `controller.MessageAction` |

---

## 📜 许可协议

This project is open-sourced under the [MIT License](LICENSE).
