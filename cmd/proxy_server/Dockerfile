FROM golang:1.22

WORKDIR /app

COPY ./go.mod ./go.sum ./RouteDefinitions.yml ./
RUN go mod download && go mod verify

COPY ./cmd/proxy_server/*.go ./cmd/proxy_server/
COPY ./internal ./internal

RUN go build -o /app/cmd/proxy_server ./cmd/proxy_server

RUN chmod +x /app/cmd/proxy_server

ENTRYPOINT ["/app/cmd/proxy_server/proxy_server"]
