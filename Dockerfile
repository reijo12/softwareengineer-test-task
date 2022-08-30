FROM golang:1.19.0-bullseye as basego

WORKDIR /scoring

COPY . .

RUN go mod download

COPY *.go ./

WORKDIR /scoring/server
RUN go build -o /docker-scoring-service

EXPOSE 9000

CMD [ "/docker-scoring-service" ]