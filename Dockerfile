FROM registry.query.consul:5000/golang:1.13.4-alpine AS builder
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o app .

FROM registry.query.consul:5000/alpine:3.9
COPY --from=builder /app/app /usr/local/bin
RUN chmod a+x /usr/local/bin/app
ENTRYPOINT [ "app" ]
