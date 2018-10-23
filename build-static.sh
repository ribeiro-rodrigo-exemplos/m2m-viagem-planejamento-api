#!/bin/sh
#Alpine image needs this static linking complile
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o m2m-viagem-planejamento-api.bin cmd/viagemPlanejamentoAPI/main.go