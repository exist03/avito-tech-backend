FROM golang:alpine

WORKDIR /avito-tech-backend/

COPY . .

RUN go mod download
EXPOSE 8080

RUN go build cmd/app/main.go

CMD ["./main"]