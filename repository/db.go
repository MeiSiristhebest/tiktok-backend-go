package repository

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/MeiSiristhebest/tiktok-backend-go/model"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	dbPath := filepath.Join(".", "tiktok.db")
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: dbLogger,
	})

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Database connection established successfully (Pure-Go SQLite/GORM)")

	// Auto Migrate Entities
	err = DB.AutoMigrate(
		&model.UserEntity{},
		&model.VideoEntity{},
		&model.FavoriteEntity{},
		&model.CommentEntity{},
		&model.RelationEntity{},
		&model.MessageEntity{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate database tables: %v", err)
	}

	seedInitialData()
}

func seedInitialData() {
	var userCount int64
	DB.Model(&model.UserEntity{}).Count(&userCount)
	if userCount > 0 {
		return
	}

	log.Println("Seeding initial demo users and videos...")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	demoUser := model.UserEntity{
		Username:        "byte_dance",
		Password:        string(hashedPassword),
		Name:            "字节跳动官方极简抖音",
		Avatar:          "https://p3-passport.byteimg.com/img/user-avatar/7f42841bc6c1e33c94ed6ec8c75d40a0~120x120.awebp",
		BackgroundImage: "https://p3-passport.byteimg.com/img/user-avatar/7f42841bc6c1e33c94ed6ec8c75d40a0~120x120.awebp",
		Signature:       "探索无界，记录美好生活！字节跳动青训营后端 Demo 服务端。",
	}
	DB.Create(&demoUser)

	demoUser2 := model.UserEntity{
		Username:        "nefelibata",
		Password:        string(hashedPassword),
		Name:            "Nefelibata (梅天翔)",
		Avatar:          "https://avatars.githubusercontent.com/u/10100000?v=4",
		BackgroundImage: "https://avatars.githubusercontent.com/u/10100000?v=4",
		Signature:       "Go Developer & AI Backend Engineer",
	}
	DB.Create(&demoUser2)

	// Seed videos (Sample MP4 URLs)
	videos := []model.VideoEntity{
		{
			AuthorID:  demoUser.ID,
			PlayURL:   "https://www.w3schools.com/html/mov_bbb.mp4",
			CoverURL:  "https://images.unsplash.com/photo-1518709268805-4e9042af9f23?w=500&auto=format&fit=crop&q=60",
			Title:     "【青训营Demo】极简版抖音 Feed 流视频体验 01",
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			AuthorID:  demoUser2.ID,
			PlayURL:   "https://www.w3schools.com/html/movie.mp4",
			CoverURL:  "https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?w=500&auto=format&fit=crop&q=60",
			Title:     "【字节Golang】基于 Gin + GORM + MySQL 打造高性能后端",
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
	}

	for _, v := range videos {
		DB.Create(&v)
	}

	log.Println("Initial demo data seeded successfully!")
}
