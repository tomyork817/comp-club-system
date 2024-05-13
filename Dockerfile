FROM golang:1.22

WORKDIR /app

COPY go.mod ./
RUN go mod download

ADD cmd ./cmd
ADD internal ./internal

RUN go build -o task ./cmd

ENTRYPOINT ["./task"]