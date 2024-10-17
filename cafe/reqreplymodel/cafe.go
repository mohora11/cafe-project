package reqreplymodel

// UserNtryReqInfo : 회원가입요청정보
type UserNtryReqInfo struct {
	CphoneNo string `json:"mobile_no" binding:"required"`   // 핸드폰번호
	LgnPwd   string `json:"login_paswd" binding:"required"` // 로그인비밀번호
}

// UserLgnReqInfo : 사용자로그인요청정보
type UserLgnReqInfo struct {
	CphoneNo string `json:"mobile_no" binding:"required"`   // 핸드폰번호
	LgnPwd   string `json:"login_paswd" binding:"required"` // 로그인비밀번호
}

// UserLgnReply : 사용자로그인회신
type UserLgnReply struct {
	AppSesTkn string `json:"app_session_token"` // 앱세션토큰
}

// UserLgotReqInfo : 사용자로그아웃요청정보
type UserLgotReqInfo struct {
	CphoneNo string `json:"mobile_no" binding:"required"` // 핸드폰번호
}

// MenuRegReqInfo : 메뉴등록요청정보
type MenuRegReqInfo struct {
	CphoneNo string `json:"mobile_no" binding:"required"`   // 핸드폰번호
	Category string `json:"category" binding:"required"`    // 카테고리
	Price    int64  `json:"price" binding:"required"`       // 가격
	Cost     int64  `json:"cost" binding:"required"`        // 원가
	MenuNm   string `json:"menu_name" binding:"required"`   // 이름
	Destn    string `json:"desctn" binding:"required"`      // 설명                     // 바코드
	Exprdt   string `json:"expiry_date" binding:"required"` // 유통기한
	Size     string `json:"size" binding:"required"`        // 사이즈
}

// MenuUptReqInfo : 메뉴수정요청정보
type MenuUptReqInfo struct {
	CphoneNo  string `json:"mobile_no" binding:"required"`    // 핸드폰번호
	MenuRegId string `json:"mn_reg_id" binding:"required"`    // 메뉴등록ID
	Category  string `json:"category" binding:"omitempty"`    // 카테고리
	Price     int64  `json:"price" binding:"omitempty"`       // 가격
	Cost      int64  `json:"cost" binding:"omitempty"`        // 원가
	MenuNm    string `json:"menu_name" binding:"omitempty"`   // 이름
	Destn     string `json:"desctn" binding:"omitempty"`      // 설명                     // 바코드
	Exprdt    string `json:"expiry_date" binding:"omitempty"` // 유통기한
	Size      string `json:"size" binding:"omitempty"`        // 사이즈
}

// MenuDelReqInfo : 메뉴삭제요청정보
type MenuDelReqInfo struct {
	CphoneNo  string `json:"mobile_no" binding:"required"` // 핸드폰번호
	MenuRegId string `json:"mn_reg_id" binding:"required"` // 메뉴등록ID
}

// MenuListInqReqInfo : 메뉴목록조회요청정보
type MenuListInqReqInfo struct {
	SrchText   string `json:"srch_text" binding:"omitempty"`  // 검색어
	CphoneNo   string `json:"mobile_no" binding:"required"`   // 핸드폰번호
	Category   string `json:"category" binding:"omitempty"`   // 카테고리
	Size       string `json:"size" binding:"required"`        // 사이즈
	PageNo     int64  `json:"page_number" binding:"required"` //페이지번호
	Pagerecnum int64  `json:"pagerecnum" binding:"required"`  //페이지당출력레코드수
}

// Menulist : 메뉴목록
type Menulist struct {
	MenuRegId string `json:"menu_reg_id" gorm:"primary_key; column:menu_reg_id"` // 메뉴등록ID
	Price     int64  `json:"price" gorm:"column:price"`                          // 가격
	Cost      int64  `json:"cost" gorm:"column:cost"`                            // 원가
	MenuNm    string `json:"menu_name" gorm:"column:menu_nm"`                    // 이름
	Size      string `json:"size" gorm:"column:size"`                            // 사이즈
}

// MenuListInqReply : 메뉴목록조회회신정보
type MenuListInqReply struct {
	MenuList []Menulist `json:"menu_list"`
}

// MenuDtlInqReqInfo : 메뉴상세조회요청정보
type MenuDtlInqReqInfo struct {
	CphoneNo  string `json:"mobile_no" binding:"required"`   // 핸드폰번호
	MenuRegId string `json:"menu_reg_id" binding:"required"` // 메뉴등록ID
}

// MenuDtlInqReply : 메뉴상세조회회신정보
type MenuDtlInqReply struct {
	MenuRegId string `json:"menu_reg_id" gorm:"column:menu_reg_id"` // 메뉴등록ID
	Category  string `json:"category" gorm:"column:category"`       // 카테고리
	Price     int64  `json:"price" gorm:"column:price"`             // 가격
	Cost      int64  `json:"cost" gorm:"column:cost"`               // 원가
	MenuNm    string `json:"menu_name" gorm:"column:menu_nm"`       // 메뉴명
	Destn     string `json:"desctn" gorm:"column:destn"`            // 설명
	Barcode   string `json:"barcode" gorm:"column:barcode"`         // 바코드
	Exprdt    string `json:"expiry_date" gorm:"column:exprdt"`      // 유통기한
	Size      string `json:"size" gorm:"column:size"`               // 사이즈
	Regdt     string `json:"registdate" gorm:"column:regdt"`        // 수정일
	Moddt     string `json:"moddate" gorm:"column:moddt"`
}
