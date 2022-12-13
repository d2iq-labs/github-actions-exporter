FROM golang:1.19 AS builder

RUN mkdir /app
WORKDIR /app
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o gha-exporter-bin .


FROM alpine:latest AS production
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/gha-exporter-bin ./
CMD ["./gha-exporter-bin"]