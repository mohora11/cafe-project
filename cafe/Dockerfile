# 빌드 환경 (경량화)
FROM golang:alpine AS builder

WORKDIR /app

# 모듈과 소스 복사
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# 애플리케이션 빌드
RUN go build -o cafeapp .

# 실행 환경 (경량화)
FROM alpine

WORKDIR /app

# 빌드된 바이너리 복사
COPY --from=builder /app/cafeapp /app/cafeapp

# 8080 포트 사용
EXPOSE 8080

# 애플리케이션 실행
CMD ["/app/cafeapp"]