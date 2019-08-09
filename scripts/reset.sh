#!/bin/sh

set -x
set -e

while true; do

  ssd init local --chain-id sschain

  ssd add-genesis-account $(sscli keys show jack -a) 1000rune,100000000stake
  ssd add-genesis-account $(sscli keys show alice -a) 1000rune,100000000stake

  sscli config chain-id sschain
  sscli config output json
  sscli config indent true
  sscli config trust-node true

  echo "password" | ssd gentx --name jack
  ssd collect-gentxs

  # add jack as a trusted account
  cat ~/.ssd/config/genesis.json | jq ".app_state.swapservice.trust_accounts[0] = {\"name\":\"Jack\", \"address\": \"$(sscli keys show jack -a)\"}" > /tmp/genesis.json
  mv /tmp/genesis.json ~/.ssd/config/genesis.json

  ssd validate-genesis

  break

done
