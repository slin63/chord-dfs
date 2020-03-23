# Launches the client
FROM golang:alpine
RUN apk add --no-cache git
ENV GODIR="/go/src/github.com/slin63" \
    CONFIG="/go/src/github.com/slin63/chord-dfs/config/config.dfs.json"

ADD . $GODIR/chord-dfs

RUN go get github.com/slin63/raft-consensus/pkg/responses
RUN go build -o dfs $GODIR/chord-dfs/cmd/client/main.go
CMD ["/bin/sh", "-c", "./dfs"]


# # Launches the entire ecosystem
# FROM golang:alpine
# RUN apk add --no-cache git
# ENV GODIR="/go/src/github.com/slin63" \
#     CONFIG="/go/src/github.com/slin63/chord-dfs/config/config.dfs.json"

# ADD . $GODIR/chord-dfs

# # RUN go get github.com/slin63/chord-failure-detector/cmd/fd
# RUN go get github.com/slin63/raft-consensus/pkg/responses
# RUN go build -o dfs $GODIR/chord-dfs/cmd/client/main.go
# CMD ["/bin/sh", "-c", "./dfs"]

