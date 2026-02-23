# Step 1: Modules caching
FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/app/main.go

# Step 3: Final Image
FROM alpine:latest
COPY --from=builder /bin/app /app
COPY --from=builder /app/config /config
COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/.env /.env

CMD ["/app"]