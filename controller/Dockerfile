FROM golang:1.22-bullseye as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /controller ./cmd/controller/main.go

FROM scratch

COPY --from=builder controller /bin/controller

CMD ["/bin/controller"]