FROM golang:1.12.10

WORKDIR /opt/e2e

COPY ./e2e/go.mod go.sum

RUN go mod download

COPY ./e2e/. .

CMD ["go", "test", "-v", "./..."]