FROM golang:1.18.2-alpine3.15 AS builder

WORKDIR /go/src/github.com/mfzl/tpl

COPY go.mod go.mod

ENV GO111MODULE on
ENV CGO_ENABLED 0

RUN go mod download

COPY . .

RUN go build -o /usr/bin/tpl

FROM scratch

COPY --from=builder /usr/bin/tpl /usr/bin/tpl

ENTRYPOINT ["/usr/bin/tpl"]
