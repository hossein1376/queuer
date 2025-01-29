FROM golang:1.23.5 AS BUILDER
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /build/querier ./cmd/querier && chmod +x /build/querier

FROM alpine:3.21.2
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR app
COPY --from=builder /build/querier /app/querier
CMD ["./querier"]