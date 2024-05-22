FROM golang:1.22

WORKDIR /app

COPY ./go.mod ./go.sum ./RouteDefinitions.yml ./
RUN go mod download && go mod verify

COPY ./cmd/token_bucket_refiller/*.go ./cmd/token_bucket_refiller/
COPY ./internal ./internal

RUN go build -o /app/cmd/token_bucket_refiller ./cmd/token_bucket_refiller

RUN chmod +x /app/cmd/token_bucket_refiller

ENTRYPOINT ["/app/cmd/token_bucket_refiller/token_bucket_refiller"]
