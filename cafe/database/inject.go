package database

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DB 미들웨어: DB 객체를 gin.Context에 담음
func InjectDB(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// DB 객체를 컨텍스트에 주입("db")
		c.Set("db", db)
		c.Next()
	}
}
