FROM golang:1.15-buster AS builder

# GO ENV VARS
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# COPY SRC
WORKDIR /build
COPY ./src .

# BUILD
RUN go build -o main .

FROM ubuntu as prod
COPY --from=builder /build/main /
CMD ["/main"]

FROM builder as test
RUN go test ./.../ -v -timeout 15m
