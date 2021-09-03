FROM golang:1.16.6-alpine3.14

WORKDIR /app

COPY  . /app

RUN go mod download

RUN go build -o server ./src/kom.com/server/rest/restGorilla/server_http.go 

EXPOSE 8080

CMD ["./server"]