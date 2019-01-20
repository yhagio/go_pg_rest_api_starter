FROM golang:1.11

WORKDIR /go/src/go_rest_pg_starter
COPY . .

RUN go get
RUN go install

CMD ["go",  "run",  "main.go"]