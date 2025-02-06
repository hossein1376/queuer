FROM golang:1.23.5 AS BUILDER
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /build/queuer ./cmd/queuer && chmod +x /build/queuer

FROM alpine:3.21.2
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR app
COPY --from=builder /build/queuer /app/queuer
CMD ["./queuer"]