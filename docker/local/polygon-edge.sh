#!/bin/sh

set -e

POLYGON_EDGE_BIN=./polygon-edge
GENESIS_BOOT_BOOTNODE_REPLACER_BIN=/polygon-edge/genesis-bootnode-replacer
GENESIS_PATH=/genesis/genesis.json
CHAIN_CUSTOM_OPTIONS=$(tr "\n" " " << EOL
--block-gas-limit 1000000000000
--epoch-size 10
--chain-id 37
--name Seedcoin
--premine 0xefb99bd12AB243Ee95BD3B8023dbAc769f56BB7A:137438953472000000000000000000
EOL
)

case "$1" in

   "init")
      case "$2" in 
         "ibft")
          if [ -f "$GENESIS_PATH" ]; then
            echo "Secrets have already been generated."
          else
            echo "Generating secrets..."
            secrets=$("$POLYGON_EDGE_BIN" secrets init --insecure --num 4 --data-dir /data/data- --json)
            echo "Secrets have been successfully generated"
            echo "Generating IBFT Genesis file..."
            cp .env /data/.env
            cd /data && /polygon-edge/polygon-edge genesis  $CHAIN_CUSTOM_OPTIONS \
              --dir /data/bootnode-source.json \
              --consensus ibft \
              --ibft-validators-prefix-path data- \
              --validator-set-size=4 \
              --bootnode /dns4/node-1/tcp/1478/p2p/$(echo $secrets | jq -r '.[0] | .node_id') \
              --bootnode /dns4/node-2/tcp/1478/p2p/$(echo $secrets | jq -r '.[1] | .node_id')
            echo "Genesis file has been successfully generated"
            echo "Copying templatefile"
            cp /data/seedcoin-template.json /data/genesis.json
            echo "Replacing bootnodes"
            go run /polygon-edge/genesis-bootnode-replacer/main.go "/data/genesis.json" "/data/bootnode-source.json"
            echo "Bootnodes replaced"
          fi
          ;;
          "polybft")
              if [ -f "$GENESIS_PATH" ]; then
                echo "Secrets have already been generated."
              else
                echo "Generating PolyBFT secrets..."
                secrets=$("$POLYGON_EDGE_BIN" polybft-secrets init --insecure --num 4 --data-dir /data/data- --json)
                echo "Secrets have been successfully generated"

                echo "Generating manifest..."
                "$POLYGON_EDGE_BIN" manifest --path /data/manifest.json --validators-path /data --validators-prefix data-

                echo "Generating PolyBFT Genesis file..."
                "$POLYGON_EDGE_BIN" genesis $CHAIN_CUSTOM_OPTIONS \
                  --dir /data/bootnode-source.json \
                  --consensus polybft \
                  --manifest /data/manifest.json \
                  --validator-set-size=4 \
                  --bootnode /dns4/node-1/tcp/1478/p2p/$(echo $secrets | jq -r '.[0] | .node_id') \
                  --bootnode /dns4/node-2/tcp/1478/p2p/$(echo $secrets | jq -r '.[1] | .node_id')
                echo "Genesis file has been successfully generated"
                echo "Copying templatefile"
                cp /data/seedcoin-template.json /data/genesis.json
                echo "Replacing bootnodes"
                go run /polygon-edge/genesis-bootnode-replacer/main.go "/genesis/genesis.json" "/genesis/bootnode-source.json"
                echo "Bootnodes replaced"
              fi
              ;;
              esac
      ;;
   *)
      echo "Executing polygon-edge..."
      exec "$POLYGON_EDGE_BIN" "$@"
      ;;

esac
