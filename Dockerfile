FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /app/segmentation ./cmd/segmentation

EXPOSE 8080

CMD ["/app/segmentation"]