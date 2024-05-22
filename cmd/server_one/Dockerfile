FROM golang:1.22

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY ./cmd/server_one/*.go ./cmd/server_one/

RUN go build -o /app/cmd/server_one ./cmd/server_one

RUN chmod +x /app/cmd/server_one

ENTRYPOINT ["/app/cmd/server_one/server_one"]
