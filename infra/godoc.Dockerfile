FROM golang:1.20-alpine

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install golang.org/x/tools/cmd/godoc@latest

ENV GOPATH /go
EXPOSE 6060
CMD ["godoc", "-http", ":6060"]
