FROM golang:1.22-alpine as builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o run ./

FROM alpine as runner

WORKDIR /app
COPY --from=builder /build/run /app/run
ENTRYPOINT [ "./run" ]