# use official Golang image
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .
# install dependencies
# RUN go get -d -v ./...

EXPOSE 8080

# RUN go build -o api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api .
RUN chmod +x ./api

# CMD ["./api"]

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/api .
COPY --from=builder /app/.env .
EXPOSE 8080

CMD ["./api"]