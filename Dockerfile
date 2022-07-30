FROM golang:1.18-alpine

EXPOSE 8080

RUN apk add make

COPY . /graffiti-berlin-svc
WORKDIR /graffiti-berlin-svc

RUN make compile

CMD ["./main"]