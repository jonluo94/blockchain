#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

LANGUAGE="golang"
CHAINCODE_VERSION="v0"
CC_SRC_PATH1="dairyfarm"
CHAINCODE_NAME1="dairyfarm"
CC_SRC_PATH2="machining"
CHAINCODE_NAME2="machining"
CC_SRC_PATH3="salesterminal"
CHAINCODE_NAME3="salesterminal"

jq --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
	echo
	exit 1
fi

echo "POST request Enroll on Org1  ..."
echo
ORG1_TOKEN=$(curl -s -X POST \
  http://localhost:8000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Jim&orgName=Org1')
echo $ORG1_TOKEN
ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "ORG1 token is $ORG1_TOKEN"
echo

echo "POST request Enroll on Org2  ..."
echo
ORG2_TOKEN=$(curl -s -X POST \
  http://localhost:8000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Jim&orgName=Org2')
echo $ORG2_TOKEN
ORG2_TOKEN=$(echo $ORG2_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "ORG2 token is $ORG2_TOKEN"
echo

echo "POST request Enroll on Org3  ..."
echo
ORG3_TOKEN=$(curl -s -X POST \
  http://localhost:8000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=Jim&orgName=Org3')
echo $ORG3_TOKEN
ORG3_TOKEN=$(echo $ORG3_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo "ORG3 token is $ORG3_TOKEN"
echo


echo "POST request Create channel mychannel ..."
echo
curl -s -X POST \
  http://localhost:8000/channels \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"channelName":"mychannel",
	"channelConfigPath":"../../network/channel/artifacts/mychannel.tx"
}'
echo
echo
sleep 5

echo "POST request Join channel on Org1"
echo
curl -s -X POST \
  http://localhost:8000/channels/mychannel/peers \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org1.example.com","peer1.org1.example.com"]
}'
echo
echo

echo "POST request Join channel on Org2"
echo
curl -s -X POST \
  http://localhost:8000/channels/mychannel/peers \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org2.example.com","peer1.org2.example.com"]
}'
echo
echo


echo "POST request Join channel on Org3"
echo
curl -s -X POST \
  http://localhost:8000/channels/mychannel/peers \
  -H "authorization: Bearer $ORG3_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org3.example.com","peer1.org3.example.com"]
}'
echo
echo


echo "POST Install chaincode on Org1"
echo
curl -s -X POST \
  http://localhost:8000/chaincodes \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.org1.example.com\",\"peer1.org1.example.com\"],
	\"chaincodeName\":\"$CHAINCODE_NAME1\",
	\"chaincodePath\":\"$CC_SRC_PATH1\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"$CHAINCODE_VERSION\"
}"
echo
echo

echo "POST Install chaincode on Org1"
echo
curl -s -X POST \
  http://localhost:8000/chaincodes \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.org1.example.com\",\"peer1.org1.example.com\"],
	\"chaincodeName\":\"$CHAINCODE_NAME2\",
	\"chaincodePath\":\"$CC_SRC_PATH2\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"$CHAINCODE_VERSION\"
}"
echo
echo

echo "POST Install chaincode on Org1"
echo
curl -s -X POST \
  http://localhost:8000/chaincodes \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.org1.example.com\",\"peer1.org1.example.com\"],
	\"chaincodeName\":\"$CHAINCODE_NAME3\",
	\"chaincodePath\":\"$CC_SRC_PATH3\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"$CHAINCODE_VERSION\"
}"
echo
echo


echo "POST Install chaincode on Org2"
echo
curl -s -X POST \
  http://localhost:8000/chaincodes \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.org2.example.com\",\"peer1.org2.example.com\"],
	\"chaincodeName\":\"$CHAINCODE_NAME1\",
	\"chaincodePath\":\"$CC_SRC_PATH1\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"$CHAINCODE_VERSION\"
}"
echo
echo

echo "POST Install chaincode on Org2"
echo
curl -s -X POST \
  http://localhost:8000/chaincodes \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.org2.example.com\",\"peer1.org2.example.com\"],
	\"chaincodeName\":\"$CHAINCODE_NAME2\",
	\"chaincodePath\":\"$CC_SRC_PATH2\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"$CHAINCODE_VERSION\"
}"
echo
echo

echo "POST Install chaincode on Org2"
echo
curl -s -X POST \
  http://localhost:8000/chaincodes \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.org2.example.com\",\"peer1.org2.example.com\"],
	\"chaincodeName\":\"$CHAINCODE_NAME3\",
	\"chaincodePath\":\"$CC_SRC_PATH3\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"$CHAINCODE_VERSION\"
}"
echo
echo


echo "POST Install chaincode on Org3"
echo
curl -s -X POST \
  http://localhost:8000/chaincodes \
  -H "authorization: Bearer $ORG3_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.org3.example.com\",\"peer1.org3.example.com\"],
	\"chaincodeName\":\"$CHAINCODE_NAME1\",
	\"chaincodePath\":\"$CC_SRC_PATH1\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"$CHAINCODE_VERSION\"
}"
echo
echo

echo "POST Install chaincode on Org3"
echo
curl -s -X POST \
  http://localhost:8000/chaincodes \
  -H "authorization: Bearer $ORG3_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.org3.example.com\",\"peer1.org3.example.com\"],
	\"chaincodeName\":\"$CHAINCODE_NAME2\",
	\"chaincodePath\":\"$CC_SRC_PATH2\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"$CHAINCODE_VERSION\"
}"
echo
echo

echo "POST Install chaincode on Org3"
echo
curl -s -X POST \
  http://localhost:8000/chaincodes \
  -H "authorization: Bearer $ORG3_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.org3.example.com\",\"peer1.org3.example.com\"],
	\"chaincodeName\":\"$CHAINCODE_NAME3\",
	\"chaincodePath\":\"$CC_SRC_PATH3\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"$CHAINCODE_VERSION\"
}"
echo
echo

echo "POST instantiate chaincode on peer1 of Org1"
echo
curl -s -X POST \
  http://localhost:8000/channels/mychannel/chaincodes \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"chaincodeName\":\"$CHAINCODE_NAME1\",
	\"chaincodeVersion\":\"$CHAINCODE_VERSION\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"args\":[\"\"]
}"
echo
echo

echo "POST instantiate chaincode on peer1 of Org2"
echo
curl -s -X POST \
  http://localhost:8000/channels/mychannel/chaincodes \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"chaincodeName\":\"$CHAINCODE_NAME2\",
	\"chaincodeVersion\":\"$CHAINCODE_VERSION\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"args\":[\"\"]
}"
echo
echo

echo "POST instantiate chaincode on peer1 of Org3"
echo
curl -s -X POST \
  http://localhost:8000/channels/mychannel/chaincodes \
  -H "authorization: Bearer $ORG3_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"chaincodeName\":\"$CHAINCODE_NAME3\",
	\"chaincodeVersion\":\"$CHAINCODE_VERSION\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"args\":[\"\"]
}"
echo
echo



echo "success ！！！"
