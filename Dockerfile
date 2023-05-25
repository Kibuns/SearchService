FROM golang:1.16-alpine

WORKDIR /app

ENV GO111MODULE=on
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY . ./

RUN go build -o /SearchService

EXPOSE 10000

CMD [ "/SearchService" ]
