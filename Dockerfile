FROM golang:1.20-alpine3.18 as builder
RUN apk --no-cache add tzdata
RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/var/cache/apt go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

FROM scratch as production
COPY --from=builder /app/main .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Europe/Moscow
CMD ["/main"]