package api

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	dbml "cafe/database/dbmodels"
	utility "cafe/utility"

	"net/http"
)

// CafeAppAPIApplyRoutes : app api applies router to the gin Engine
func CafeAppAPIApplyRoutes(r *gin.RouterGroup) {

	// 회원가입과 로그인에는 미들웨어를 적용하지 않음
	r.POST("/userntry", UserNtryReq) // 사용자가입요청
	r.POST("/login", UserLoginReq)   // 로그인요청

	// 나머지 요청에는 Validator 미들웨어 적용
	cafeRoutes := r.Group("")   // /cafe 하위 라우팅
	cafeRoutes.Use(Validator()) // 이 라우팅 그룹에 미들웨어 적용

	// 다른 API들에는 미들웨어가 적용됨
	cafeRoutes.POST("/logout", UserLogoutReq)       // 로그아웃요청
	cafeRoutes.POST("/menureg", MenuRegReq)         // 메뉴등록요청
	cafeRoutes.POST("/menudel", MenuDelReq)         // 메뉴삭제요청
	cafeRoutes.POST("/menuupt", MenuUptReq)         // 메뉴수정요청
	cafeRoutes.POST("/menulistinq", MenuListInqReq) // 메뉴목록조회요청
	cafeRoutes.POST("/menudtlinq", MenuDtlInqReq)   // 메뉴상세조회요청

}

// 미들웨어로 토큰 검증 적용
func Validator() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Gin Context에서 DB를 가져옴
		db, exists := c.Get("db")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"meta": gin.H{
					"code":    http.StatusInternalServerError,
					"message": "DB연결에 실패하였습니다. 잠시 후 다시 시도해주세요.",
				},
			})
			c.Abort()
			return
		}

		token := c.GetHeader("t") // 헤더에서 토큰을 가져옴
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"meta": gin.H{
					"code":    http.StatusUnauthorized,
					"message": "앱세션토큰을 읽어오는데 실패하였습니다. 잠시 후 다시 시도해주세요.",
				},
			})
			c.Abort()
			return
		}

		isValid, errMsg := TokenCheck(db.(*gorm.DB), token)
		if !isValid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"meta": gin.H{
					"code":    http.StatusUnauthorized,
					"message": errMsg,
				},
			})
			c.Abort()
			return
		}
	}
}

// TokenCheck : 토큰검증
func TokenCheck(db *gorm.DB, tokenSource string) (isValid bool, errMsg string) {

	var (
		jwtKey   = []byte("1021503")
		UserInfo dbml.UserInfo
	)

	// 토큰을 파싱하여 만료 시간을 확인
	claims := &utility.Claims{}
	_, err := jwt.ParseWithClaims(tokenSource, claims, func(token *jwt.Token) (interface{}, error) {

		// 서명 키 반환
		return jwtKey, nil
	})

	if err != nil {
		// 만료된 토큰
		return false, "세션이 만료되었습니다. 다시 로그인 해주시기 바랍니다."
	}

	//휴대폰번호가 빈 값일 경우
	if claims.CphoneNo == "" {
		return false, "권한이 없습니다."
	}

	// 사용자가 존재하는지 확인
	if err := db.Table("user_info").Where("cphone_no = ?", claims.CphoneNo).Find(&UserInfo).Error; err != nil {

		return false, "권한이 없습니다."
	}

	// 토큰이 유효한 경우
	return true, ""
}
