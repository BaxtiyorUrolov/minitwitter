# Go image'ni asos sifatida olish
FROM golang:1.22.1 AS builder

# Ishchi katalogni yaratish
WORKDIR /app

# Modullarni va kodni yuklash
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Dasturiy ta'minotni yaratish
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o main ./cmd/main.go

# Yana birinchi rasmdan alohida rasm yaratish
FROM alpine:latest

# Ishga tushirish uchun portni ochish
EXPOSE 7070

# Asosiy faylni ko'chirish
WORKDIR /root/
COPY --from=builder /app .

# Dasturiy ta'minotni ishga tushirish
CMD ["./main"]
