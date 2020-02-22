FROM golang:alpine
ENV CONFIG="config.json"
ENV INTRODUCER=false

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main src/main.go
RUN go build -o membership chord-failure-detector/src/main.go

## TODO: have this start up the remote logger as well
CMD ["sh", "-c", "./scripts/init.sh"]
