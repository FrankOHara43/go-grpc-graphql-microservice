FROM golang:1.25-alpine AS build

RUN apk --no-cache add gcc g++ make ca-certificates

WORKDIR /go/src/github.com/FrankOHara43/go-grpc-graphql-microservice

COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
COPY catalog catalog
COPY order order

RUN GO111MODULE=on go build -mod=vendor -o /go/bin/app ./order/cmd/order

FROM alpine:3.19

WORKDIR /usr/bin

COPY --from=build /go/bin/app .

EXPOSE 8080

CMD ["app"]