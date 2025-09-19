FROM golang:1.25

WORKDIR /app

RUN apt-get update && apt-get install -y binutils && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build .

ENTRYPOINT ["/app/elfGetArg"]