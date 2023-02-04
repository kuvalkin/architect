FROM golang:1.19-bullseye as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /app/executable .

RUN chmod +x /app/executable

FROM debian:bullseye

WORKDIR /app

COPY --from=builder /app/executable .

EXPOSE 8000

CMD ["./executable"]