# Launches the entire ecosystem

FROM golang:alpine
RUN apk add --no-cache git
ENV INTRODUCER=0 \
    GODIR="/go/src/github.com/slin63" \
    ALL=1

RUN mkdir -p /temp/store/

# ADD to prevent docker from caching out of date git repos
ADD https://api.github.com/repos/slin63/chord-failure-detector/git/refs/heads/dfs version.json
RUN git clone -b dfs https://github.com/slin63/chord-failure-detector $GODIR/chord-failure-detector/

ADD https://api.github.com/repos/slin63/raft-consensus/git/refs/heads/master version.json
RUN git clone -b master https://github.com/slin63/raft-consensus $GODIR/raft-consensus/

ADD . $GODIR/chord-dfs

RUN go build -o member $GODIR/chord-failure-detector/cmd/fd/main.go
RUN go build -o dfs $GODIR/chord-dfs/cmd/server/main.go
RUN go build -o raft $GODIR/raft-consensus/cmd/raft/main.go

CMD ["sh", "-c", "$GODIR/chord-dfs/scripts/init.sh"]

