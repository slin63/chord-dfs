FROM golang:alpine
RUN apk add --no-cache git
ENV CONFIG="/go/src/github.com/slin63/chord-dfs/config.json" INTRODUCER=0

RUN mkdir -p /temp/store/
WORKDIR /go/src/github.com/slin63/
RUN git clone https://github.com/slin63/chord-failure-detector

ADD . /go/src/github.com/slin63/chord-dfs

RUN go build -o dfs ./chord-dfs/cmd/dfs/main.go
RUN go build -o member ./chord-failure-detector/cmd/fd/main.go

CMD ["sh", "-c", "./chord-dfs/scripts/init.sh"]

