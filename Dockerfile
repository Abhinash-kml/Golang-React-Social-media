
ARG VERSION=1.23.1
FROM golang:${VERSION}
WORKDIR /app
COPY . .
RUN go mod download
WORKDIR /app/cmd
EXPOSE 8000/tcp
CMD ["go", "run", "main.go"]