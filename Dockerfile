FROM golang:1.20.5-alpine3.18 AS builder

WORKDIR /log-parser

RUN apk add --no-cache git make

COPY go.mod .

RUN go mod download

COPY . .

RUN go mod vendor

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o /go/bin/log-parser .

ENTRYPOINT ["sh", "init.sh"]

FROM alpine:3.18

COPY --from=builder /go/bin/log-parser /usr/bin/log-parser
COPY --from=builder /log-parser/resources /etc/log-parser/resources/
COPY --from=builder /log-parser/init.sh /usr/bin/init.sh

ENTRYPOINT ["sh", "/usr/bin/init.sh"]