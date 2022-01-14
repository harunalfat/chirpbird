FROM golang:1.17-alpine as build

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . .
RUN go build -o /app/binary .

FROM alpine
COPY --from=build /app/binary /app/binary

CMD ["/app/binary"]