FROM golang:1.18-alpine3.17 AS builder

RUN go version
ENV GOPATH=/

COPY . /github.com/zh0vtyj/allincecup-server/
WORKDIR /github.com/zh0vtyj/allincecup-server/

RUN go mod download
RUN go mod download github.com/ugorji/go
RUN GOOS=linux go build -o alliancecup ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/zh0vtyj/allincecup-server/alliancecup .

CMD ["./alliancecup"]