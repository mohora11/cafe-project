package main

import (
	"cafe/api"
	"cafe/database"

	"log"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {

	// DB 연결 host DB와 연결을 위해 host.docker.internal 사용
	dsn := "root:z1s2c3f4##@tcp(host.docker.internal)/cafe_db?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Gin 기본 웹 서버
	ginEngine := gin.Default()

	// DB 미들웨어 등록
	ginEngine.Use(database.InjectDB(DB)) // DB를 gin.Context에 주입

	// API 라우팅 설정
	api.ApplyRoutes(ginEngine)

	// 서버 실행
	ginEngine.Run(":8080")
}
