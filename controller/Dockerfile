FROM golang:1.22-bullseye as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o controller ./cmd/controller/main.go

FROM scratch

ARG ENVIRONMENT
ENV CONFIG_PATH=/usr/local/etc/${ENVIRONMENT}.yml

COPY --from=builder /app/controller /usr/local/bin/controller
COPY --from=builder /app/configs/${ENVIRONMENT}.yml ${CONFIG_PATH}

CMD ["/usr/local/bin/controller"]