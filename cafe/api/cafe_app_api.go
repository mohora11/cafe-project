package api

import (
	dbml "cafe/database/dbmodels"
	model "cafe/reqreplymodel"
	utility "cafe/utility"
	"log"

	"time"

	"database/sql"

	"gorm.io/gorm"
)

// UserNtry : 회원등록
func UserNtry(db *gorm.DB, pRequest *model.UserNtryReqInfo) (bool, string) {

	if len(pRequest.LgnPwd) > 0 && len(pRequest.CphoneNo) > 0 {

		// 핸드폰 번호 중복 체크
		var existUser dbml.UserInfo
		if err := db.Where("cphone_no = ?", pRequest.CphoneNo).First(&existUser).Error; err == nil {

			return false, "이미 등록된 핸드폰 번호입니다."
		}

		lgnPwdHash := utility.MakeHash(pRequest.LgnPwd)

		// 회원등록
		srvrTime, _ := utility.GetSrvrTime(db)
		NtryUser := dbml.UserInfo{
			CphoneNo: pRequest.CphoneNo,
			LgnPwd:   lgnPwdHash,
			Ntrydt:   srvrTime,
			Moddt:    sql.NullTime{},
		}
		if err := db.Create(&NtryUser).Error; err != nil {
			return false, "회원등록에 실패하였습니다. 잠시 후 다시 시도해주세요."
		}

	} else {
		return false, "핸드폰 번호나 비밀번호를 다시 입력해주세요."
	}

	return true, "회원등록이 완료되었습니다. 로그인 화면에서 로그인 해주세요."
}

// UserLogin : 회원 로그인
func UserLogin(db *gorm.DB, pRequest *model.UserLgnReqInfo) (bool, string, model.UserLgnReply) {

	var (
		UserInfo dbml.UserInfo

		resp model.UserLgnReply
	)

	srvrTime, _ := utility.GetSrvrTime(db)

	// 요청정보 체크
	if len(pRequest.LgnPwd) == 0 && len(pRequest.CphoneNo) == 0 {
		return false, "핸드폰 번호나 비밀번호를 다시 입력해주세요.", model.UserLgnReply{}
	}
	lgnPwdHash := utility.MakeHash(pRequest.LgnPwd)

	d1 := db.Select("cphone_no, lgn_pwd, ntrydt, moddt").Table("user_info").Where("cphone_no = ?", pRequest.CphoneNo)
	if err := d1.Find(&UserInfo).Error; err != nil {
		return false, "일치하는 회원정보가 없습니다.", resp
	}

	// Hash처리한 요청 비밀번호, db의 Hash처리된 비밀번호 비교
	if UserInfo.LgnPwd == lgnPwdHash {
		token, err := utility.MakeToken(db, pRequest.CphoneNo, pRequest.LgnPwd)
		if err != nil {
			return false, "세션토큰 생성 실패", resp
		}
		// 생성된 토큰을 회신 구조체에 할당
		resp = model.UserLgnReply{
			AppSesTkn: token,
		}

		// 로그인할때쓴 세션토큰을 DB에 저장 - 로그아웃시 삭제를 위함
		err = db.Table(`user_info`).Where(`cphone_no = ?`, pRequest.CphoneNo).Updates(map[string]interface{}{
			"ses_tkn": token,    // 세션토큰
			"moddt":   srvrTime, // 수정일
		}).Error

		if err != nil {
			return false, "로그인 오류 발생 잠시 후 다시 시도해 주세요.", resp
		}

	} else {
		return false, "비밀번호가 일치하지 않습니다", resp
	}

	return true, "로그인 성공", resp
}

// UserLogout : 회원 로그아웃
func UserLoout(db *gorm.DB, token string, pRequest *model.UserLgotReqInfo) (bool, string) {

	var (
		dbSestkn string
	)

	srvrTime, _ := utility.GetSrvrTime(db)

	// 요청정보 체크
	if len(pRequest.CphoneNo) == 0 {
		return false, "로그아웃에 실패하였습니다 잠시 후 다시 시도해 주세요."
	}

	// 기존 세션토큰 조회
	d1 := db.Select("ses_tkn").Table(`user_info`).Where(`cphone_no = ?`, pRequest.CphoneNo)
	if err := d1.Take(&dbSestkn).Error; err != nil {
		return false, "일치하는 회원정보가 없습니다."
	}

	// 프로그램 강제종료로 인한 db토큰이 삭제되지 않았을경우 이전에 저장되어있던 토큰과 비교
	if dbSestkn != token {
		return false, "이미 로그아웃 상태입니다. 다시 로그인 해주세요."
	}

	// 저장되어있던 세션토큰 삭제
	err := db.Table(`user_info`).Where(`cphone_no = ?`, pRequest.CphoneNo).Updates(map[string]interface{}{
		"ses_tkn": "",       // 세션토큰 삭제
		"moddt":   srvrTime, // 수정일
	}).Error

	if err != nil {
		return false, "로그아웃 실패 잠시 후 다시 시도해주세요."
	}

	return true, "로그아웃 성공"
}

// MenuReg : 메뉴등록
func MenuReg(db *gorm.DB, token string, pRequest *model.MenuRegReqInfo) (bool, string) {

	var (
		dbSestkn string
	)

	if len(pRequest.Category) == 0 && len(pRequest.MenuNm) == 0 && len(pRequest.Destn) == 0 && len(pRequest.Size) == 0 {
		return false, "메뉴등록정보를 다시 입력해 주세요."
	}

	// 유통기한 파싱
	expiryDate, err := time.Parse(utility.TimeFormatDay, pRequest.Exprdt)
	if err != nil {
		return false, "유통기한 입력 형식이 잘못되었습니다."
	}

	// 기존 세션토큰 조회
	d1 := db.Select("ses_tkn").Table(`user_info`).Where(`cphone_no = ?`, pRequest.CphoneNo)
	if err := d1.Take(&dbSestkn).Error; err != nil {
		return false, "일치하는 회원정보가 없습니다."
	}

	// 프로그램 강제종료로 인한 db토큰이 삭제되지 않았을경우 새 로그인 세션 토큰이랑 비교
	if dbSestkn != token {
		return false, "세션정보가 만료 되었습니다 다시 로그인 해주세요."
	}

	// 메뉴 등록 ID 생성
	mnRegId := utility.MakeBizID("Mn", 6)

	log.Println(pRequest.MenuNm)

	// 메뉴 초성명 생성
	menuChosung := utility.ExtractChoSung(pRequest.MenuNm)

	// 바코드 생성
	newbrcd := utility.GenerateBarcode(pRequest.Category, pRequest.MenuNm, pRequest.Price, mnRegId)

	// 메뉴 등록
	srvrTime, _ := utility.GetSrvrTime(db)
	MenuRegi := dbml.MenuInfo{
		MenuRegId: mnRegId,
		Category:  pRequest.Category,
		Price:     pRequest.Cost,
		Cost:      pRequest.Cost,
		MenuNm:    pRequest.MenuNm,
		MenuNmCs:  menuChosung,
		Destn:     pRequest.Destn,
		Barcode:   newbrcd,
		Exprdt:    expiryDate,
		Size:      pRequest.Size,
		Regdt:     srvrTime,
		Moddt:     sql.NullTime{},
	}
	if err := db.Create(&MenuRegi).Error; err != nil {
		return false, "메뉴 등록에 실패하였습니다 잠시 후 다시 시도해주세요."
	}

	return true, "메뉴 등록이 완료 되었습니다."
}

// MenuUpt : 메뉴수정
func MenuUpt(db *gorm.DB, token string, pRequest *model.MenuUptReqInfo) (bool, string) {

	var (
		dbSestkn string

		newExpiryDate time.Time
		err           error

		dbMenuInfo dbml.MenuInfo
	)

	// 기존 세션토큰 조회
	d1 := db.Select("ses_tkn").Table(`user_info`).Where(`cphone_no = ?`, pRequest.CphoneNo)
	if err := d1.Take(&dbSestkn).Error; err != nil {
		return false, "일치하는 회원정보가 없습니다."
	}

	// 프로그램 강제종료로 인한 db토큰이 삭제되지 않았을경우 새 로그인 세션 토큰이랑 비교
	if dbSestkn != token {
		return false, "기존 회원정보와 일치하지 않습니다 다시 로그인 해주세요."
	}

	// 기존 메뉴 정보 조회

	if err := db.Table("menu_info").Where("menu_reg_id = ?", pRequest.MenuRegId).First(&dbMenuInfo).Error; err != nil {

		return false, "기존 메뉴 정보를 찾지 못했습니다. 잠시 후 다시 시도 해주세요."
	}

	// 수정할 값들 설정 (요청값이 있는 경우에만 업데이트)
	if pRequest.MenuNm != "" {
		dbMenuInfo.MenuNm = pRequest.MenuNm

		//초성명도 수정
		dbMenuInfo.MenuNmCs = utility.ExtractChoSung(pRequest.MenuNm)
	}
	if pRequest.Category != "" {
		dbMenuInfo.Category = pRequest.Category
	}
	if pRequest.Price != 0 {
		dbMenuInfo.Price = pRequest.Price
	}
	if pRequest.Cost != 0 {
		dbMenuInfo.Cost = pRequest.Cost
	}
	if pRequest.Destn != "" {
		dbMenuInfo.Destn = pRequest.Destn
	}
	if pRequest.Exprdt != "" {

		// 유통기한 파싱
		newExpiryDate, err = time.Parse(utility.TimeFormatDay, pRequest.Exprdt)
		if err != nil {
			return false, "유통기한 입력 형식이 잘못되었습니다."
		}

		dbMenuInfo.Exprdt = newExpiryDate
	}
	if pRequest.Size != "" {
		dbMenuInfo.Size = pRequest.Size
	}

	// 신규 바코드 생성
	newbrcd := utility.GenerateBarcode(dbMenuInfo.Category, dbMenuInfo.MenuNm, dbMenuInfo.Price, pRequest.MenuRegId)

	// 현재서버시간
	srvrTime, _ := utility.GetSrvrTime(db)

	// 메뉴 수정 (필드 변경 후 업데이트)
	d2 := db.Table("menu_info").Where("menu_reg_id = ?", pRequest.MenuRegId).Updates(map[string]interface{}{
		"category":   dbMenuInfo.Category,
		"price":      dbMenuInfo.Price,
		"cost":       dbMenuInfo.Cost,
		"menu_nm":    dbMenuInfo.MenuNm,
		"menu_cs_nm": dbMenuInfo.MenuNmCs,
		"destn":      dbMenuInfo.Destn,
		"barcode":    newbrcd,
		"exprdt":     dbMenuInfo.Exprdt,
		"size":       dbMenuInfo.Size,
		"moddt":      srvrTime,
	})

	if err := d2.Error; err != nil {
		return false, "메뉴 업데이트에 실패하였습니다. 잠시 후 다시 시도해 주세요."
	}

	return true, "메뉴 업데이트가 완료 되었습니다."
}

// MenuDel : 메뉴삭제
func MenuDel(db *gorm.DB, token string, pRequest *model.MenuDelReqInfo) (bool, string) {

	var (
		dbSestkn string
	)

	// 기존 세션토큰 조회
	d1 := db.Select("ses_tkn").Table(`user_info`).Where(`cphone_no = ?`, pRequest.CphoneNo)
	if err := d1.Take(&dbSestkn).Error; err != nil {
		return false, "일치하는 회원정보가 없습니다."
	}

	// 프로그램 강제종료로 인한 db토큰이 삭제되지 않았을경우 새 로그인 세션 토큰이랑 비교
	if dbSestkn != token {
		return false, "기존 회원정보와 일치하지 않습니다 다시 로그인 해주세요."
	}

	// 메뉴 삭제
	err := db.Where("menu_reg_id = ?", pRequest.MenuRegId).Delete(dbml.MenuInfo{}).Error
	if err != nil {
		return false, "메뉴 삭제에 실패하였습니다. 잠시 후 다시 시도해 주세요."
	}

	return true, "메뉴 삭제가 완료 되었습니다."
}

// MenuListInq : 메뉴목록조회
func MenuListInq(db *gorm.DB, token string, pRequest *model.MenuListInqReqInfo) (bool, string, model.MenuListInqReply) {

	var (
		dbSestkn string

		resp model.MenuListInqReply
	)

	log.Println(pRequest)

	// 기존 세션토큰 조회
	d1 := db.Select("ses_tkn").Table(`user_info`).Where(`cphone_no = ?`, pRequest.CphoneNo)
	if err := d1.Take(&dbSestkn).Error; err != nil {
		return false, "일치하는 회원정보가 없습니다.", resp
	}

	// 프로그램 강제종료로 인한 db토큰이 삭제되지 않았을경우 새 로그인 세션 토큰이랑 비교
	if dbSestkn != token {
		return false, "기존 회원정보와 일치하지 않습니다 다시 로그인 해주세요.", resp
	}

	// 메뉴 목록조회(페이징)
	d2 := db.Scopes(utility.Paginate(int(pRequest.PageNo), int(pRequest.Pagerecnum), "price ASC")).
		Select("menu_reg_id, price, cost, menu_nm, size").Table("menu_info").Where("size = ?", pRequest.Size)

	// 카테고리 요청값
	if len(pRequest.Category) > 0 {
		d2 = d2.Where("category = ?", pRequest.Category)
	}

	// 검색어 요청값
	if len(pRequest.SrchText) > 0 {

		if utility.IsChoSung(pRequest.SrchText) {
			// 초성 검색
			d2 = d2.Where(`menu_cs_nm LIKE ?`, "%"+pRequest.SrchText+"%")
		} else {
			// 본이름 검색
			d2 = d2.Where(`menu_nm LIKE ?`, "%"+pRequest.SrchText+"%")
		}
	}

	if err := d2.Scan(&resp.MenuList).Error; err != nil {

		return false, "메뉴 목록 검색 중 오류가 발생했습니다. 잠시 후 다시 시도 해주세요.", resp
	}

	return true, "ok", resp
}

// MenuDtlInq : 메뉴상세조회
func MenuDtlInq(db *gorm.DB, token string, pRequest *model.MenuDtlInqReqInfo) (bool, string, model.MenuDtlInqReply) {

	var (
		dbSestkn string

		resp       model.MenuDtlInqReply
		dbMenuInfo dbml.MenuInfo
	)

	// 기존 세션토큰 조회
	d1 := db.Select("ses_tkn").Table(`user_info`).Where(`cphone_no = ?`, pRequest.CphoneNo)
	if err := d1.Take(&dbSestkn).Error; err != nil {
		return false, "일치하는 회원정보가 없습니다.", resp
	}

	// 프로그램 강제종료로 인한 db토큰이 삭제되지 않았을경우 새 로그인 세션 토큰이랑 비교
	if dbSestkn != token {
		return false, "기존 회원정보와 일치하지 않습니다 다시 로그인 해주세요.", resp
	}

	// 메뉴상세정보조회
	d2 := db.Table("menu_info").Where("menu_reg_id = ?", pRequest.MenuRegId)

	err := d2.Take(&dbMenuInfo).Error
	if err != nil {
		return false, "상품의 상세정보를 가져오는 중 오류가 발생했습니다. 잠시 후 다시 시도해주세요.", resp
	}

	resp.MenuRegId = dbMenuInfo.MenuRegId
	resp.Category = dbMenuInfo.Category
	resp.Price = dbMenuInfo.Price
	resp.Cost = dbMenuInfo.Cost
	resp.MenuNm = dbMenuInfo.MenuNm
	resp.Destn = dbMenuInfo.Destn
	resp.Barcode = dbMenuInfo.Barcode
	resp.Exprdt = dbMenuInfo.Exprdt.Format(utility.TimeFormatDay)
	resp.Size = dbMenuInfo.Size
	resp.Regdt = dbMenuInfo.Regdt.Format(utility.TimeFormatDay)

	// sql.Nulltime 처리(수정한적이없다면 null값)
	if dbMenuInfo.Moddt.Valid {
		resp.Moddt = dbMenuInfo.Moddt.Time.Format(utility.TimeFormatDay)
	} else {
		resp.Moddt = ""
	}

	return true, "ok", resp
}
