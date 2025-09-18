#!/bin/bash
echo -e "CONTAINER \t ID NAME \t SERVICE" && for container_id in $(docker ps -q); do service=$(docker inspect --format '{{ index .Config.Labels "com.docker.compose.service" }}' $container_id); name=$(docker inspect --format '{{ .Name }}' $container_id | sed 's/\///g'); echo -e "$container_id \t $name \t $service"; done
