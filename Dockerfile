FROM golang:1.22

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /receipt-processor

EXPOSE 8080

ENTRYPOINT ["/receipt-processor"]
