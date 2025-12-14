FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /todo-app-final

ENV TODO_PORT=7540
ENV TODO_PASSWORD=12345

EXPOSE 7540

CMD ["/todo-app-final"]
