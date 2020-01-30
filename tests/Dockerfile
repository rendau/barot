FROM golang:1.13 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go test -c -o tests



FROM debian:stable

WORKDIR /app

COPY --from=builder /app/tests ./
COPY gherkin ./gherkin

CMD ["./tests"]
