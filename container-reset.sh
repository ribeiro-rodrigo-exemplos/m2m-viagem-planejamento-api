#!/bin/sh
docker container stop m2m-viagem-planejamento-api 
docker container rm m2m-viagem-planejamento-api 
docker container run -d -p 8081:8081 --name m2m-viagem-planejamento-api m2m/m2m-viagem-planejamento-api:nightly