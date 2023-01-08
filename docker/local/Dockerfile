FROM golang:1.18-alpine AS builder

WORKDIR /polygon-edge

ADD go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o polygon-edge main.go

FROM golang:1.18-alpine AS runner

RUN apk --no-cache add ca-certificates jq

WORKDIR /polygon-edge

COPY ./genesis-bootnode-replacer/main.go /polygon-edge/genesis-bootnode-replacer/main.go
COPY --from=builder /polygon-edge/polygon-edge ./
COPY ./docker/local/polygon-edge.sh ./
COPY ./chain/chains/seedcoin.json /genesis/seedcoin-template.json
RUN chmod 0777 /genesis/seedcoin-template.json
RUN chmod 0777 polygon-edge.sh

# Expose json-rpc, libp2p and grpc ports
EXPOSE 8545 9632 1478 5001

ENTRYPOINT ["./polygon-edge.sh"]
