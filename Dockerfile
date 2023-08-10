# syntax=docker/dockerfile:1
FROM golang:latest

WORKDIR /usr/src/app

RUN apt update
RUN apt install -y nodejs npm

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/todo-app .
RUN cd frontend && npm install
RUN cd frontend && npm run build

EXPOSE 80

CMD [ "todo-app" ]