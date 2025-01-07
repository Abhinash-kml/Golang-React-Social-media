FROM golang:1.23.1
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd
RUN go build -o /app/my-app
EXPOSE 8080
CMD ["/app/my-app"]