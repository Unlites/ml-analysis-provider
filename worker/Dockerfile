FROM golang:1.22-bullseye as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o worker ./cmd/worker/main.go

FROM scratch

ARG ENVIRONMENT
ENV CONFIG_PATH=/usr/local/etc/${ENVIRONMENT}.yml

COPY --from=builder /app/worker /usr/local/bin/worker
COPY --from=builder /app/configs/${ENVIRONMENT}.yml ${CONFIG_PATH}

CMD ["/usr/local/bin/worker"]