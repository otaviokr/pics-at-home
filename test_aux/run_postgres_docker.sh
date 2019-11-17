#!/bin/bash

# $HOME/docker/volumes/postgres
DB_VOLUME=$1

INSTANCE_NAME="pg-docker"
POSTGRES_PASSWORD="docker"
POSTGRES_PORT=5432

if [[ ${DB_VOLUME} == "" ]]; then
    echo "No volume defined."
else 
    echo "Volume defined. Database data will be stored at: ${DB_VOLUME}"
fi

docker run -d \
        --rm \
        --name ${INSTANCE_NAME} \
        -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
        -p $POSTGRES_PORT:5432 \
        -v ${DB_VOLUME}:/var/lib/postgresql/data \
        postgres