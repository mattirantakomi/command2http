FROM golang:1.16-alpine as builder

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" .

FROM scratch
COPY --from=builder /src/command2http /
