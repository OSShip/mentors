FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY utils/ /app/utils/
COPY services/mentors/ /app/services/mentors/
WORKDIR /app/services/mentors
RUN go mod download && CGO_ENABLED=0 go build -o /mentors .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget
COPY --from=builder /mentors /mentors
EXPOSE 8085
CMD ["/mentors"]
