package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	model "cafe/reqreplymodel"
)

// UserNtryReq : 사용자 가입 요청
func UserNtryReq(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 요청 정보
	reqInfo := model.UserNtryReqInfo{}

	// json 요청정보 바인딩
	if err := c.BindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": gin.H{
				"code":    http.StatusBadRequest,
				"message": "잘못된 요청입니다. 요청 형식을 확인해주세요.",
			},
			"data": nil,
		})
		return
	}

	// 회원가입 서비스
	resultFlag, errMsg := UserNtry(db, &reqInfo)

	// 정상 처리 및 오류 정보 회신
	if resultFlag {

		c.JSON(http.StatusOK, gin.H{
			"meta": gin.H{
				"code":    http.StatusOK,
				"message": errMsg,
			},
		})
	} else {

		c.JSON(http.StatusBadRequest, gin.H{
			"meta": gin.H{
				"code":    http.StatusBadRequest,
				"message": errMsg,
			},
		})
	}
}

// UserLoginReq : 사용자 로그인 요청
func UserLoginReq(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 요청 정보
	reqInfo := model.UserLgnReqInfo{}

	// json 요청정보 바인딩
	if err := c.BindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": gin.H{
				"code":    http.StatusBadRequest,
				"message": "잘못된 요청입니다. 요청 형식을 확인해주세요.",
			},
		})
		return
	}

	// 로그인 서비스
	resultFlag, errMsg, resp := UserLogin(db, &reqInfo)

	// 정상 처리 및 오류 정보 회신
	if resultFlag {

		c.JSON(http.StatusOK, gin.H{
			"meta": gin.H{
				"code":    http.StatusOK,
				"message": errMsg,
			},
			"data": gin.H{
				"message": resp,
			},
		})
	} else {

		c.JSON(http.StatusInternalServerError, gin.H{
			"meta": gin.H{
				"code":    http.StatusInternalServerError,
				"message": errMsg,
			},
			"data": resp,
		})
	}
}

// UserLogoutReq : 사용자 로그아웃 요청
func UserLogoutReq(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 요청 정보
	reqInfo := model.UserLgotReqInfo{}

	// json 요청정보 바인딩
	if err := c.BindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": gin.H{
				"code":    http.StatusBadRequest,
				"message": "잘못된 요청입니다. 요청 형식을 확인해주세요.",
			},
		})
		return
	}

	// 헤더토큰 가져오기
	tokenKey := c.Request.Header.Get("t")

	resultFlag, errMsg := UserLoout(db, tokenKey, &reqInfo)

	// 정상 처리 및 오류 정보 회신
	if resultFlag {

		c.JSON(http.StatusOK, gin.H{
			"meta": gin.H{
				"code":    http.StatusOK,
				"message": errMsg,
			},
		})
	} else {

		c.JSON(http.StatusBadRequest, gin.H{
			"meta": gin.H{
				"code":    http.StatusBadRequest,
				"message": errMsg,
			},
		})
	}
}

// MenuRegReq : 메뉴등록요청
func MenuRegReq(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 요청 정보
	reqInfo := model.MenuRegReqInfo{}

	// json 요청정보 바인딩
	if err := c.BindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": gin.H{
				"code":    http.StatusBadRequest,
				"message": "잘못된 요청입니다. 요청 형식을 확인해주세요.",
			},
		})
		return
	}

	// 헤더토큰 가져오기
	tokenKey := c.Request.Header.Get("t")

	// 메뉴등록 서비스
	resultFlag, errMsg := MenuReg(db, tokenKey, &reqInfo)

	// 정상 처리 및 오류 정보 회신
	if resultFlag {

		c.JSON(http.StatusOK, gin.H{
			"meta": gin.H{
				"code":    http.StatusOK,
				"message": errMsg,
			},
		})
	} else {

		c.JSON(http.StatusInternalServerError, gin.H{
			"meta": gin.H{
				"code":    http.StatusInternalServerError,
				"message": errMsg,
			},
		})
	}
}

// MenuUptReq : 메뉴수정요청
func MenuUptReq(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 요청 정보
	reqInfo := model.MenuUptReqInfo{}

	// json 요청정보 바인딩
	if err := c.BindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": gin.H{
				"code":    http.StatusBadRequest,
				"message": "잘못된 요청입니다. 요청 형식을 확인해주세요.",
			},
		})
		return
	}

	// 헤더토큰 가져오기
	tokenKey := c.Request.Header.Get("t")

	// 메뉴업데이트 서비스
	resultFlag, errMsg := MenuUpt(db, tokenKey, &reqInfo)

	// 정상 처리 및 오류 정보 회신
	if resultFlag {

		c.JSON(http.StatusOK, gin.H{
			"meta": gin.H{
				"code":    http.StatusOK,
				"message": errMsg,
			},
		})
	} else {

		c.JSON(http.StatusInternalServerError, gin.H{
			"meta": gin.H{
				"code":    http.StatusInternalServerError,
				"message": errMsg,
			},
		})
	}
}

// MenuDelReq : 메뉴삭제요청
func MenuDelReq(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 요청 정보
	reqInfo := model.MenuDelReqInfo{}

	// json 요청정보 바인딩
	if err := c.BindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": gin.H{
				"code":    http.StatusBadRequest,
				"message": "잘못된 요청입니다. 요청 형식을 확인해주세요.",
			},
		})
		return
	}

	// 헤더토큰 가져오기
	tokenKey := c.Request.Header.Get("t")

	// 메뉴삭제 서비스
	resultFlag, errMsg := MenuDel(db, tokenKey, &reqInfo)

	// 정상 처리 및 오류 정보 회신
	if resultFlag {

		c.JSON(http.StatusOK, gin.H{
			"meta": gin.H{
				"code":    http.StatusOK,
				"message": errMsg,
			},
		})
	} else {

		c.JSON(http.StatusInternalServerError, gin.H{
			"meta": gin.H{
				"code":    http.StatusInternalServerError,
				"message": errMsg,
			},
		})
	}
}

// MenuListInqReq : 메뉴목록조회요청
func MenuListInqReq(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 요청 정보
	reqInfo := model.MenuListInqReqInfo{}

	// json 요청정보 바인딩
	if err := c.BindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": gin.H{
				"code":    http.StatusBadRequest,
				"message": "잘못된 요청입니다. 요청 형식을 확인해주세요.",
			},
		})
		return
	}

	// 헤더토큰 가져오기
	tokenKey := c.Request.Header.Get("t")

	// 메뉴목록조회 서비스
	resultFlag, errMsg, resp := MenuListInq(db, tokenKey, &reqInfo)

	// 정상 처리 및 오류 정보 회신
	if resultFlag {

		c.JSON(http.StatusOK, gin.H{
			"meta": gin.H{
				"code":    http.StatusOK,
				"message": errMsg,
			},
			"data": gin.H{
				"products": resp,
			},
		})
	} else {

		c.JSON(http.StatusInternalServerError, gin.H{
			"meta": gin.H{
				"code":    http.StatusInternalServerError,
				"message": errMsg,
			},
			"data": gin.H{
				"products": resp,
			},
		})
	}
}

// MenuDtlInqReq : 메뉴상세조회요청
func MenuDtlInqReq(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// 요청 정보
	reqInfo := model.MenuDtlInqReqInfo{}

	// json 요청정보 바인딩
	if err := c.BindJSON(&reqInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"meta": gin.H{
				"code":    http.StatusBadRequest,
				"message": "잘못된 요청입니다. 요청 형식을 확인해주세요.",
			},
		})
		return
	}

	// 헤더토큰 가져오기
	tokenKey := c.Request.Header.Get("t")

	// 메뉴상세조회 서비스
	resultFlag, errMsg, resp := MenuDtlInq(db, tokenKey, &reqInfo)

	// 정상 처리 및 오류 정보 회신
	if resultFlag {

		c.JSON(http.StatusOK, gin.H{
			"meta": gin.H{
				"code":    http.StatusOK,
				"message": errMsg,
			},
			"data": gin.H{
				"products": resp,
			},
		})
	} else {

		c.JSON(http.StatusInternalServerError, gin.H{
			"meta": gin.H{
				"code":    http.StatusInternalServerError,
				"message": errMsg,
			},
			"data": gin.H{
				"products": resp,
			},
		})
	}
}
