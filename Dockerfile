FROM golang:1.14-alpine as go

WORKDIR /src/go
COPY go/go.mod .
COPY go/go.sum .
RUN go mod download

COPY go/. .
RUN go build -o server main.go

FROM node:14 as node

WORKDIR /src/ui
COPY ui/package.json .
COPY ui/yarn.lock .
RUN yarn install

COPY ui/. .
RUN yarn build

FROM alpine

WORKDIR /workdir
COPY --from=go /src/go/server server
COPY --from=node /src/ui/build/. /workdir/public/.

CMD [ "./server" ]
