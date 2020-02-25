FROM golang:alpine
ENV CONFIG="config.json" INTRODUCER=0 SERVER=0

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main src/main.go

## TODO: have this start up the remote logger as well
CMD ["sh", "-c", "./scripts/init.sh"]
