FROM golang:1.25-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o account ./cmd

EXPOSE 8001

CMD ["/app/account"]
