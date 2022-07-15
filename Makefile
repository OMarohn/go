PROJECT?=kom.com/m/v2/server_http
APP?=src/kom.com/server/rest/echo/server_http.go
APPBIN?=server_http
RELEASE?=1.0.4
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
BUILD_USER?=$(shell whoami)
PORT?=8080
LDFLAGS?="-s -w -X main.version=${RELEASE} -X main.commit=${COMMIT} -X main.date=${BUILD_TIME} -X main.builtBy=${BUILD_USER}"

build: 
	go build -ldflags ${LDFLAGS} -o bin/${APPBIN} ${APP}

run: build
	PORT=${PORT} ./bin/${APPBIN}

runtest:
	go test -v -race ./...

docker: 
	docker build --build-arg ldflags=${LDFLAGS} -t omarohn/coaster-server:v${RELEASE} -f Dockerfile.multistage . 	
	docker push omarohn/coaster-server:v${RELEASE}

