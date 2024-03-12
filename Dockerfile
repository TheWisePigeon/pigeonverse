FROM golang:latest as builder
RUN update-ca-certificates
WORKDIR app/
COPY go.mod .
ENV GO111MODULE=on
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -extldflags "-static"' -o /app .
FROM alpine:latest
WORKDIR /usr/local/bin
COPY --from=builder /app /usr/local/bin/app
COPY views /usr/local/bin/views
RUN mkdir assets
COPY assets/dracula.min.css /usr/local/bin/assets/dracula.css
COPY assets/output.css /usr/local/bin/assets/output.css
COPY assets/styles.css /usr/local/bin/assets/styles.css
COPY posts /usr/local/bin/posts
ENV CONTENT_DIR=/usr/local/bin/posts/
ENV PORT=1000
ENV ENV=PROD
EXPOSE 1000
ENTRYPOINT ["./app"]
