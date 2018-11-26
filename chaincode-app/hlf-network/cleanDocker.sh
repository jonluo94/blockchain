#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.
set -e

# Shut down the Docker containers for the system tests.
docker-compose down

# remove chaincode docker images
docker rm -f $(docker ps -aq)
docker rmi $(docker images *-hellocc-* -q)

# Your system is now clean