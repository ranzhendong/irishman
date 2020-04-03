FROM golang:1.13.8-alpine

ADD . /home/irishman

WORKDIR /home/irishman

RUN go mod tidy

RUN go build -v -o /home/irishman/irishman

CMD ["./irishman"]