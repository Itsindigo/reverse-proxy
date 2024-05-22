FROM golang:1.22

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

COPY ./cmd/server_two/*.go ./cmd/server_two/

RUN go build -o /app/cmd/server_two ./cmd/server_two

RUN chmod +x /app/cmd/server_two

ENTRYPOINT ["/app/cmd/server_two/server_two"]
