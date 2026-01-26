FROM golang:1.25-alpine

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o risk-analysis ./cmd

EXPOSE 8002

CMD [ "/app/risk-analysis" ]
