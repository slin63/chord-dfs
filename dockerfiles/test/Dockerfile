# Only launches the DFS and Failure Detector

FROM golang:alpine
RUN apk add --no-cache git
ENV CONFIG="/go/src/github.com/slin63/chord-dfs/config.json" INTRODUCER=0 GODIR="/go/src/github.com/slin63" ALL=0

RUN mkdir -p /temp/store/
RUN go get github.com/slin63/raft-consensus/pkg/responses

RUN git clone https://github.com/slin63/chord-failure-detector $GODIR/chord-failure-detector/

ADD . $GODIR/chord-dfs

RUN go build -o dfs $GODIR/chord-dfs/cmd/server/main.go
RUN go build -o member $GODIR/chord-failure-detector/cmd/fd/main.go

CMD ["sh", "-c", "$GODIR/chord-dfs/scripts/init.sh"]

