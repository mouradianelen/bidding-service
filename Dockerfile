FROM golang:1.20
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o /docker-gs-ping
EXPOSE 8080
CMD ["/docker-gs-ping"]