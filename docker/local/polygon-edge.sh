#!/bin/sh

set -e

POLYGON_EDGE_BIN=./polygon-edge
GENESIS_BOOT_BOOTNODE_REPLACER_BIN=/polygon-edge/genesis-bootnode-replacer

case "$1" in

   "init")
      echo "Generating secrets..."
      secrets=$("$POLYGON_EDGE_BIN" secrets init --num 4 --data-dir data- --json)
      echo "Secrets have been successfully generated"

      echo "Generating genesis file..."
      "$POLYGON_EDGE_BIN" genesis \
        --dir /genesis/bootnode-source.json \
        --consensus ibft \
        --name Seedcoin \
        --chain-id 37 \
        --premine 0xefb99bd12AB243Ee95BD3B8023dbAc769f56BB7A:1000000000000000000000 \
        --block-gas-limit 1000000000 \
        --epoch-size 50 \
        --ibft-validators-prefix-path data- \
        --bootnode /dns4/node-1/tcp/1478/p2p/$(echo $secrets | jq -r '.[0] | .node_id') \
        --bootnode /dns4/node-2/tcp/1478/p2p/$(echo $secrets | jq -r '.[1] | .node_id')
      echo "Genesis file has been successfully generated"
      echo "Copying templatefile"
      cp /genesis/seedcoin-template.json /genesis/genesis.json
      echo "Replacing bootnodes"
      chmod 0777 /genesis/genesis.json
      chmod 0777 /genesis/bootnode-source.json

      go run /polygon-edge/genesis-bootnode-replacer/main.go "/genesis/genesis.json" "/genesis/bootnode-source.json"
      ;;
   *)
      echo "Executing polygon-edge..."
      exec "$POLYGON_EDGE_BIN" "$@"
      ;;

esac
