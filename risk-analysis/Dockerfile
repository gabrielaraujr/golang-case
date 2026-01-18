FROM golang:1.25-alpine

WORKDIR /app

COPY . .

RUN go build -o risk-analysis ./cmd

EXPOSE 8002

CMD ["./risk-analysis"]
