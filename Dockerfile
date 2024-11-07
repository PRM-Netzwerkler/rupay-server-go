FROM golang:1.22

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o /exec

EXPOSE 8888

CMD [ "/exec" ]