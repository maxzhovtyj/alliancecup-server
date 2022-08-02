FROM golang:1.16-alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o alliancecup ./cmd/main.go

CMD ["./alliancecup"]