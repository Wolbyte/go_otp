FROM golang:1.25.0-alpine

WORKDIR /app

COPY . .

RUN ["go", "mod", "download"]

EXPOSE 8080

CMD ["go", "run", "main.go"]