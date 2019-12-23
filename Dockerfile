FROM golang

ENV PORT=8080
EXPOSE 8080

WORKDIR /go/src/textSave

ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .
RUN go build -o /go/bin/textSave .

ENTRYPOINT ["/go/bin/textSave"]