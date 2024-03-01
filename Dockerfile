FROM golang:latest as builder
RUN update-ca-certificates
WORKDIR app/
COPY go.mod .
ENV GO111MODULE=on
RUN go mod download && go mod verify
COPY . .
RUN go build -o /app .
FROM debian:latest
COPY --from=builder /app /usr/local/bin/app
COPY views /usr/local/bin/views
COPY assets /usr/local/bin/assets
COPY posts /usr/local/bin/posts
WORKDIR /usr/local/bin
ENV CONTENT_DIR=/usr/local/bin/posts/
EXPOSE 1000
ENTRYPOINT ["app"]
