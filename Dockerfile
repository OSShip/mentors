FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY utils/ /app/utils/
COPY services/mentors/ /app/services/mentors/
WORKDIR /app/services/mentors
RUN go mod download && CGO_ENABLED=0 go build -o /mentors .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget \
    && addgroup -g 1001 -S osship \
    && adduser -u 1001 -S osship -G osship
COPY --from=builder /mentors /mentors
USER 1001
EXPOSE 8085
CMD ["/mentors"]
