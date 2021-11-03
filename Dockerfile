# Build golang project that listen on port 8080 and use tor
FROM golang:1.17.2-buster
RUN apt update && apt upgrade -y && apt install -y tor
WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /docker-gs-ping

EXPOSE 8080

ARG SCHEME
ARG HOST

CMD [ "/docker-gs-ping" ]