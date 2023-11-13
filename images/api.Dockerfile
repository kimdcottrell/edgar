FROM golang:1.21 as builder 

WORKDIR /go/src

# These are more expensive operations, 
#    so split them out so they can be cached in Docker layers
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
COPY /api ./api

RUN go build -o main .

FROM builder AS dev

COPY --from=qmcgaw/binpot:dlv-v1.21.0 /bin /go/bin/dlv
COPY --from=qmcgaw/binpot:gopls-v0.13.2 /bin /go/bin/gopls
COPY --from=qmcgaw/binpot:golangci-lint-v1.54.1 /bin /go/bin/golangci-lint

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH \
    CGO_ENABLED=0 \
    GO111MODULE=on

WORKDIR /go/src

EXPOSE 8080
CMD ["./main"]
