  
  echo "##########################################################"
  echo "##### Generate certificates using cryptogen tool #########"
  echo "##########################################################"

  if [ -d "channel/crypto-config" ]; then
    rm -Rf channel/crypto-config
  fi
  set -x
  ./bin/cryptogen generate  --output=./channel/crypto-config --config=./channel/cryptogen.yaml
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    echo "Failed to generate certificates..."
    exit 1
  fi

  echo "##########################################################"
  echo "#########  Generating Orderer Genesis block ##############"
  echo "##########################################################"
  # Note: For some unknown reason (at least for now) the block file can't be
  # named orderer.genesis.block or the orderer will fail to launch!
  if [ -d "channel/artifacts" ]; then
    rm -Rf channel/artifacts
  fi
  mkdir ./channel/artifacts
  set -x
  ./bin/configtxgen -configPath ./channel  -profile OrgsOrdererGenesis -outputBlock ./channel/artifacts/genesis.block
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    echo "Failed to generate orderer genesis block..."
    exit 1
  fi

  echo
  echo "#################################################################"
  echo "### Generating channel configuration transaction 'mychannel.tx' ###"
  echo "#################################################################"
  set -x
  ./bin/configtxgen -configPath ./channel -profile OrgsChannel -outputCreateChannelTx ./channel/artifacts/mychannel.tx -channelID mychannel
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    echo "Failed to generate channel configuration transaction..."
    exit 1
  fi

  echo
  echo "#################################################################"
  echo "#######    Generating anchor peer update for Org1MSP   ##########"
  echo "#################################################################"
  set -x
  ./bin/configtxgen -configPath ./channel -profile OrgsChannel -outputAnchorPeersUpdate ./channel/artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    echo "Failed to generate anchor peer update for Org1MSP..."
    exit 1
  fi

  echo
  echo "#################################################################"
  echo "#######    Generating anchor peer update for Org2MSP   ##########"
  echo "#################################################################"
  set -x
  ./bin/configtxgen -configPath ./channel -profile OrgsChannel -outputAnchorPeersUpdate ./channel/artifacts/Org2MSPanchors.tx -channelID mychannel -asOrg Org2MSP
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    echo "Failed to generate anchor peer update for Org2MSP..."
    exit 1
  fi

  echo
  echo "#################################################################"
  echo "#######    Generating anchor peer update for Org3MSP   ##########"
  echo "#################################################################"
  set -x
  ./bin/configtxgen -configPath ./channel -profile OrgsChannel -outputAnchorPeersUpdate ./channel/artifacts/Org3MSPanchors.tx -channelID mychannel -asOrg Org3MSP
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    echo "Failed to generate anchor peer update for Org3MSP..."
    exit 1
  fi
