#!/bin/bash

echo "### Recreating docker containers\n"

docker-compose up -d --force-recreate
