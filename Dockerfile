FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build --ldflags='-s -w' -o metadata_federation_service

EXPOSE 9889

CMD ["./metadata_federation_service"]