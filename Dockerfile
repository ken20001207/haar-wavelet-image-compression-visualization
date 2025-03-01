FROM golang:1.17.0-alpine3.14
RUN mkdir /src
WORKDIR /src
COPY . /src/
RUN go build -o ./bin/server
EXPOSE 8080
CMD ["./bin/server"]
