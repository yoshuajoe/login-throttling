FROM golang:1.14.4-buster
WORKDIR /go/src/telkom_test2
COPY . /go/src/telkom_test2
RUN env GOOS=linux GOARCH=amd64 go build -o cmd/main cmd/main.go cmd/middleware.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app/
COPY --from=0 /go/src/telkom_test2/cmd/main /app/main
RUN chmod +x /app/main && \
    mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
CMD ["./main"]