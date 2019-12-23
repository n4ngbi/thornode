# adds an account node into the genesis file
add_node_account () {
    NODE_ADDRESS=$1
    VALIDATOR=$2
    NODE_PUB_KEY=$3
    VERSION=$4
    BOND_ADDRESS=$5
    jq --arg VERSION $VERSION --arg BOND_ADDRESS "$BOND_ADDRESS" --arg VALIDATOR "$VALIDATOR" --arg NODE_ADDRESS "$NODE_ADDRESS" --arg NODE_PUB_KEY "$NODE_PUB_KEY" '.app_state.thorchain.node_accounts += [{"node_address": $NODE_ADDRESS, "version": $VERSION, "status":"active", "active_block_height": "0", "bond_address":$BOND_ADDRESS,"validator_cons_pub_key":$VALIDATOR,"pub_key_set":{"secp256k1":$NODE_PUB_KEY,"ed25519":$NODE_PUB_KEY}}]' <~/.thord/config/genesis.json >/tmp/genesis.json
    mv /tmp/genesis.json ~/.thord/config/genesis.json
}

add_last_event_id () {
    echo "Adding last event id $1"
    jq --arg ID $1 '.app_state.thorchain.last_event_id = $ID' ~/.thord/config/genesis.json > /tmp/genesis.json
    mv /tmp/genesis.json ~/.thord/config/genesis.json
}

add_admin_config () {
    jq --arg NODE_ADDRESS "$3" --arg KEY "$1" --arg VALUE "$2" '.app_state.thorchain.admin_configs += [{"address": $NODE_ADDRESS ,"key":$KEY, "value":$VALUE}]' ~/.thord/config/genesis.json > /tmp/genesis.json
}

# inits a thorchain with a comman separate list of usernames
init_chain () {
    export IFS=","

    thord init local --chain-id thorchain
    thorcli keys list

    for user in $@; do # iterate over our list of comma separated users "alice,jack"
        thord add-genesis-account $user 1000thor
    done

    thorcli config chain-id thorchain
    thorcli config output json
    thorcli config indent true
    thorcli config trust-node true
}

peer_list () {
    PEER="$1@$2:26656"
    ADDR='addr_book_strict = true'
    ADDR_STRICT_FALSE='addr_book_strict = false'
    PEERSISTENT_PEER_TARGET='persistent_peers = ""'

    sed -i -e "s/$ADDR/$ADDR_STRICT_FALSE/g" ~/.thord/config/config.toml
    sed -i -e "s/$PEERSISTENT_PEER_TARGET/persistent_peers = \"$PEER\"/g" ~/.thord/config/config.toml
}

gen_bnb_address () {
    if [ ! -f ~/.signer/private_key.txt ]; then
        echo "GENERATING BNB ADDRESSES"
        # because the generate command can get API rate limited, THORNode may need to retry
        n=0
        until [ $n -ge 60 ]; do
            generate > /tmp/bnb && break
            n=$[$n+1]
            sleep 1
        done
        ADDRESS=$(cat /tmp/bnb | grep MASTER= | awk -F= '{print $NF}')
        echo $ADDRESS > ~/.signer/address.txt
        BINANCE_PRIVATE_KEY=$(cat /tmp/bnb | grep MASTER_KEY= | awk -F= '{print $NF}')
        echo $BINANCE_PRIVATE_KEY > ~/.signer/private_key.txt
        PUBKEY=$(cat /tmp/bnb | grep MASTER_PUBKEY= | awk -F= '{print $NF}')
        echo $PUBKEY > ~/.signer/pubkey.txt
    fi
}



fetch_genesis () {
    until curl -s "$1:26657" > /dev/null; do
        sleep 3
    done
    curl -s $1:26657/genesis | jq .result.genesis > ~/.thord/config/genesis.json
    thord validate-genesis
    cat ~/.thord/config/genesis.json
}

fetch_node_id () {
    until curl -s "$1:26657" > /dev/null; do
        sleep 3
    done
    curl -s $1:26657/status | jq -r .result.node_info.id
}

fetch_version () {
    thorcli query thorchain version --chain-id thorchain --trust-node --output json | jq -r .version
}
