FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . /app
WORKDIR /app/myapp


RUN go build -o /myapp .

CMD [ "/myapp" ]
