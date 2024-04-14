FROM golang:1.22.0-alpine AS builder

COPY . /build
WORKDIR /build

RUN go mod download && go mod tidy
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o app ./cmd/app/main.go

FROM scratch AS runner

WORKDIR /api/
COPY --from=builder /build/app .

EXPOSE 8080

ENTRYPOINT ["./app"]