#!/bin/bash

NAME="Tasarruf"
REMOTE_IMAGE="177403612005.dkr.ecr.us-east-1.amazonaws.com/tasarruf-prod:latest"

echo "=== Stopping all running containers"

for container_id in $(docker ps -q); do docker stop $container_id;done

echo "--- Stopped all running containres"

echo "=== Remove all non running containers"

for container_id in $(docker ps -q -a); do docker rm $container_id;done

echo "--- Removed all non running containers"

echo "=== Pulling remote image from ECR"

docker pull $REMOTE_IMAGE

echo "=== Cleaning untagged images"

for image_id in $(docker images --filter "dangling=true" -q);do sudo docker image rm $image_id;done

echo "--- Cleaned old images tagged"

echo "=== Running the Docker image"

server=$(docker images ${REMOTE_IMAGE} -q)

echo "Running image with ID = ${server[0]}"

docker run  --network=host ${server[0]}
