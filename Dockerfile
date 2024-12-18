FROM golang:1.23 as build
WORKDIR /app
COPY . .
RUN go build -o app
ENTRYPOINT ["./app"]