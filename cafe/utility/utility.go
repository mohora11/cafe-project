package utility

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	TimeFormatDay = "2006-01-02"
)

// JWT 키
var jwtKey = []byte("1021503")

// JWT claims 구조체
type Claims struct {
	CphoneNo string `json:"cphone_no"`
	LgnPwd   string `json:"login_pwd"`
	jwt.StandardClaims
}

// JWT 토큰 생성
func MakeToken(db *gorm.DB, cPhoneNo, lgnPwd string) (string, error) {
	// 만료 시간을 설정
	srvrTime, _ := GetSrvrTime(db)
	expirationTime := srvrTime.Add(1 * time.Hour) // 10분

	claims := &Claims{
		CphoneNo: cPhoneNo,
		LgnPwd:   lgnPwd,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "cafeApp", // 발급자
		},
	}

	// JWT 토큰 생성
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Print(token)
	return token.SignedString(jwtKey)
}

// MakeHash : Authorized blocks unauthorized requestrs
func MakeHash(val string) string {
	hasher := sha256.New()
	hasher.Write([]byte(val))
	data := hasher.Sum(nil)
	retVal := hex.EncodeToString(data)
	return retVal
}

// GetSrvrTime: 서버시각 가져오기
func GetSrvrTime(db *gorm.DB) (time.Time, string) {
	type dummy struct {
		SrvrTime time.Time `json:"-" gorm:"column:srvr_time"`
	}

	var reDb dummy
	var errMessage string = ""

	row := db.Raw("select now() srvr_time")

	if err := row.Take(&reDb).Error; err != nil {
		errMessage = err.Error()
	}

	return reDb.SrvrTime, errMessage
}

// StringWithCharset :
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// MakeStrRandom : 랜덤stirng 생성
func MakeStrRandom(length int) string {
	return StringWithCharset(length, charset)
}

// 초성 테이블 정의
var cho = []rune{
	'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ', 'ㅅ',
	'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ',
}

// 초성을 추출하는 함수
func ExtractChoSung(input string) string {
	var choSungPart []rune

	// 한글의 초성만 추출
	for _, r := range input {
		// 한글 유니코드 범위 확인 (가 ~ 힣)
		if r >= 0xAC00 && r <= 0xD7A3 {
			// 초성 추출: (유니코드 - 0xAC00) / (21 * 28) 계산
			choIndex := (r - 0xAC00) / (21 * 28)
			choSungPart = append(choSungPart, cho[choIndex]) // 초성만 모은다
		}
	}

	return string(choSungPart)
}

// contains : 주어진 문자 `r`이 초성 리스트에 포함되는지 확인
func contains(slice []rune, r rune) bool {
	for _, c := range slice {
		if c == r {
			return true
		}
	}
	return false
}

// IsChoSung 초성 확인 함수
func IsChoSung(input string) bool {
	for _, r := range input {
		// 초성 유니코드 범위 확인
		if r < 0x3131 || r > 0x314E {
			// 초성이 아닌 문자가 있으면 false 반환
			return false
		}
	}
	return true
}

// MakeBizID : 메뉴 고유 등록 ID 생성
func MakeBizID(pHeaderID string, pLastIDLength int) string {
	var returnID string = ""
	var byteBuff bytes.Buffer
	byteBuff.WriteString(pHeaderID)
	byteBuff.WriteString(MakeStrRandom(pLastIDLength))
	returnID = byteBuff.String()
	byteBuff.Reset()
	return returnID
}

// 바코드 생성 함수
func GenerateBarcode(category, menuName string, price int64, menuRegId string) string {
	return fmt.Sprintf("%s-%s-%s-%s", category, menuName, strconv.FormatInt(price, 10), menuRegId)
}

// Paginate : 페이징
func Paginate(pPage int, pPageSize int, pOrderString string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pPage == 0 {
			pPage = 1
		}
		switch {
		case pPageSize > 999:
			pPageSize = 999
		case pPageSize <= 0:
			pPageSize = 10
		}
		offset := (pPage - 1) * pPageSize

		return db.Offset(offset).Limit(pPageSize).Order(pOrderString)
	}
}

// // MakeToken : 앱 세션용 토큰생성
// func MakeToken(
// 	db *gorm.DB,
// 	userNo string, //회원번호
// 	cphoneNo string, //휴대폰번호
// ) (
// 	tokenSource string, //토큰소스
// 	tokenExpireDateTimeString string, //토큰만료시간
// 	succeed bool, //토큰생성 성공여부
// 	errCd string, //에러코드
// 	errMsg string, //에러메시지
// ) {

// 	var (
// 		err            error
// 	)

// 	// 토큰타입에 따른 사전검사

// 	//회원번호
// 	if len(userNo) < 1 {
// 		return tokenSource, tokenExpireDateTimeString, false, "", err.Error()
// 	}

// 	parts := []string{
// 		userNo,        // 회원번호
// 		cphoneNo,      // 휴대폰번호
// 	}
// 	tokenSource = strings.Join(parts, "||")

// 	return tokenSource, tokenExpireDateTimeString, true, "", ""
// }
