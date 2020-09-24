FROM golang:1.12.10 as build-env

ENV CONFIG_PATH $CONFIG_PATH

WORKDIR /action-log

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /opt/action-log

# Release
FROM alpine:latest

WORKDIR /root/

COPY --from=build-env /opt/action-log .
COPY --from=build-env /action-log/config ./config

CMD ["./action-log", "run", "$CONFIG_PATH"]