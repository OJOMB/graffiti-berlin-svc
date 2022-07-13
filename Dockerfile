FROM golang:1.18-alpine

EXPOSE 8080

RUN apk add make

COPY . /subscription-svc
WORKDIR /subscription-svc

RUN make compile

CMD ["./main"]