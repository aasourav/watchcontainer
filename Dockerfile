FROM golang:1.24.1-bookworm AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/watcher/main.go

FROM scratch
COPY --from=build /app/server /server
ENTRYPOINT ["/server"]
