#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
export PATH=~/bin:$PATH
export FABRIC_CFG_PATH=${PWD}
export VERBOSE=false
export CHANNEL_NAME=mychannel
rm -rf channel-artifacts
mkdir channel-artifacts
rm -rf crypto-config
mkdir crypto-config

function generateCerts() {
  which cryptogen
  if [ "$?" -ne 0 ]; then
    echo "cryptogen tool not found. exiting"
    exit 1
  fi
  echo
  echo "##########################################################"
  echo "##### Generate certificates using cryptogen tool #########"
  echo "##########################################################"
  cryptogen generate --config=./crypto-config.yaml
}

function generateChannelArtifacts() {
  which configtxgen
  if [ "$?" -ne 0 ]; then
    echo "configtxgen tool not found. exiting"
    exit 1
  fi

  echo "##########################################################"
  echo "#########  Generating Orderer Genesis block ##############"
  echo "##########################################################"
  configtxgen -profile AllOrgsOrdererGenesis -channelID syschannel -outputBlock ./channel-artifacts/genesis.block
  
  
  echo "#################################################################"
  echo "### Generating channel configuration transaction 'channel.tx' ###"
  echo "#################################################################"
  configtxgen -profile AllOrgsChannel -outputCreateChannelTx ./channel-artifacts/mychannel.tx -channelID $CHANNEL_NAME

  echo
  echo "#################################################################"
  echo "#######    Generating anchor peer update for Org1MSP   ##########"
  echo "#################################################################"
  configtxgen -profile AllOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP

  echo
  echo "#################################################################"
  echo "#######    Generating anchor peer update for Org2MSP   ##########"
  echo "#################################################################"
  configtxgen -profile AllOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP
}
generateCerts
generateChannelArtifacts
