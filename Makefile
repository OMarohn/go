APP?=src/kom.com/server/rest/echo/server_http.go
APPBIN?=server_http
LABEL?=v1.0.4
PORT?=8080

build: 
	go build -o bin/${APPBIN} ${APP}

run: build
	PORT=${PORT} ./bin/${APP}

runtest:
	go test -v -race ./...

docker:
	docker build -t omarohn/coaster-server:${LABEL} -f Dockerfile.multistage . 	
	docker push omarohn/coaster-server:${LABEL}

