FROM golang:1.14.4-buster
WORKDIR /go/src/login-throttling
COPY . /go/src/login-throttling
RUN env GOOS=linux GOARCH=amd64 go build -o cmd/main cmd/main.go cmd/middleware.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app/
COPY --from=0 /go/src/login-throttling/cmd/main /app/main
RUN chmod +x /app/main && \
    mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
CMD ["./main"]