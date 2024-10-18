# report cafe

gin framework를 이용한 cafe 앱을 만들어 보았습니다

이전회사의 프로젝트를 떠올리며 비슷한 구조로 만들었으나
자주쓰는 함수들을 utility.go에 따로 다 담아 다소 보기 불편할 수 있겠습니다

api - cafe_api_service : 핸들러 함수
      cafe_app_api : 비즈니스 로직
      cafe_route_handler: 각 서비스 라우터와 토큰 검증 미들웨어


database - dbmodels - cafe_model : 데이터베이스 테이블을 구조체로 마이그레이션

reqreplymodel - cafe : json요청정보 및 회신정보 구조체 마이그레이션

utility - utility : 서버시간, 한국어 초성 검증, JWT토큰 생성 등 유틸성 함수 모음집
