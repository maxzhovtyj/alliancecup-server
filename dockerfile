FROM golang:1.18-alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go mod download github.com/ugorji/go
RUN go build -o alliancecup ./cmd/main.go

CMD ["./alliancecup"]