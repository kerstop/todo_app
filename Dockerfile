# syntax=docker/dockerfile:1
FROM golang:latest as development

WORKDIR /usr/src/app

EXPOSE 80

RUN apt update
RUN apt install -y nodejs npm

COPY go.mod go.sum ./
RUN go mod download && go mod verify

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM development as build

COPY . .
RUN go build -v -o /usr/local/bin/todo-app .
RUN cd frontend && npm install
RUN cd frontend && npm run build

FROM scratch as production

COPY --from=build /usr/local/bin/todo-app /usr/local/bin/todo-app
COPY --from=build /usrc/src/app/frontend/dist /usrc/src/app/frontend/dist

CMD /usr/local/bin/todo-app 