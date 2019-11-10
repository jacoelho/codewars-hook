FROM golangci/golangci-lint:v1.21.0 as lint 

WORKDIR /app 

COPY . .

RUN golangci-lint run ./...


FROM golang:1.13-stretch as builder

WORKDIR /build

RUN useradd \
    --no-log-init \
    --no-create-home \
    --shell /sbin/nologin \
    -u 1000 app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go test -v ./...

RUN go build -trimpath -ldflags "-linkmode external -extldflags -static" ./cmd/...

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /build/codewars-hook /codewars-hook

USER app

ENTRYPOINT ["/codewars-hook"]