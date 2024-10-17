package cafe

import (
	"database/sql"
	"time"
)

// UserInfo : 회원정보
type UserInfo struct {
	CphoneNo string       `json:"mobile_no" gorm:"cphone_no"`          // 핸드폰번호
	LgnPwd   string       `json:"login_paswd" gorm:"column:lgn_pwd"`   // 로그인비밀번호
	Ntrydt   time.Time    `json:"entry_date" gorm:"column:ntrydt"`     // 가입일
	Moddt    sql.NullTime `json:"moddate" gorm:"column:moddt"`         // 수정일
	SesTkn   string       `json:"session_token" gorm:"column:ses_tkn"` // 사용중인세션토큰
}

func (UserInfo) TableName() string {
	return "user_info"
}

// MenuInfo : 메뉴정보
type MenuInfo struct {
	MenuRegId string       `json:"menu_reg_id" gorm:"primary_key; column:menu_reg_id"` // 메뉴등록ID
	Category  string       `json:"category" gorm:"primary_key; column:category"`       // 카테고리
	Price     int64        `json:"price" gorm:"column:price"`                          // 가격
	Cost      int64        `json:"cost" gorm:"column:cost"`                            // 원가
	MenuNm    string       `json:"menu_name" gorm:"column:menu_nm"`                    // 메뉴명
	MenuNmCs  string       `json:"menu_chosung_name" gorm:"column:menu_cs_nm"`         // 메뉴초성명
	Destn     string       `json:"desctn" gorm:"column:destn"`                         // 설명
	Barcode   string       `json:"barcode" gorm:"column:barcode"`                      // 바코드
	Exprdt    time.Time    `json:"expiry_date" gorm:"column:exprdt"`                   // 유통기한
	Size      string       `json:"size" gorm:"column:size"`                            // 사이즈
	Regdt     time.Time    `json:"registdate" gorm:"column:regdt"`                     // 수정일
	Moddt     sql.NullTime `json:"moddate" gorm:"column:moddt"`                        // 수정일
}

func (MenuInfo) TableName() string {
	return "menu_info"
}
