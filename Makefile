CURRENT_DIR=$(shell pwd)


build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${ELD_DATA} ${CURRENT_DIR}/cmd/server/main.go

build-image:
	docker build --rm -t ${REGISTRY}/${PROJECT_NAME}/${APP}/${IMG_NAME}:${TAG} .
	docker tag ${REGISTRY}/${PROJECT_NAME}/${APP}/${IMG_NAME}:${TAG} ${REGISTRY}/${PROJECT_NAME}/${APP}/${IMG_NAME}:${TAG_LATEST}

run:
	go run cmd/main.go

swag-init:
	swag init -g api/router.go -o api/docs
