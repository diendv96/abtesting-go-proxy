FROM golang:1.13 as builder

WORKDIR /proxy

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . ./
# Building using -mod=vendor
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM alpine:3.10
RUN apk add --update ca-certificates
RUN apk add --no-cache tzdata && \
  cp -f /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
  apk del tzdata
WORKDIR /app
EXPOSE 8000
COPY ./config.yml /app/
COPY --from=builder /proxy/app .
ENTRYPOINT ["./app"]
