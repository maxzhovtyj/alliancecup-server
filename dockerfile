FROM golang:1.18-alpine3.17 AS builder

RUN go version
ENV GOPATH=/

COPY . /github.com/zh0vtyj/alliancecup-server/
WORKDIR /github.com/zh0vtyj/alliancecup-server/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/alliancecup ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/zh0vtyj/alliancecup-server/.bin/alliancecup .
COPY --from=0 /github.com/zh0vtyj/alliancecup-server/configs/config.yml configs/
COPY --from=0 /github.com/zh0vtyj/alliancecup-server/.env .

EXPOSE 3000

CMD ["./alliancecup"]