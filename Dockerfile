FROM golang:1.13 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN make



FROM debian:stable

WORKDIR /app

COPY --from=builder /app/cmd/build/* ./

EXPOSE 80

CMD ["./barot_api"]
