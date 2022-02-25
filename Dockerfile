FROM golang:1.16.13-alpine3.15

WORKDIR /app

COPY  . /app

RUN go mod download

RUN go build -o server ./src/kom.com/server/rest/restGorilla/server_http.go 

EXPOSE 8080

RUN addgroup -S gocoaster && adduser -S 12345 --uid 12345 -G gocoaster

USER 12345

CMD ["./server"]