FROM golang:1.12.10 as build-env

ENV CONFIG_PATH $CONFIG_PATH

WORKDIR /action

COPY ./api/go.mod go.sum

RUN go mod download

COPY ./api/. .

RUN CGO_ENABLED=0 go build -o /opt/action

# Release
FROM alpine:latest

WORKDIR /root/

COPY --from=build-env /opt/action .
COPY --from=build-env /action/config ./config

CMD ["./action", "run", "$CONFIG_PATH"]