FROM golang:1.22-bullseye as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /worker ./cmd/worker/main.go

FROM scratch

COPY --from=builder worker /bin/worker

CMD ["/bin/worker"]