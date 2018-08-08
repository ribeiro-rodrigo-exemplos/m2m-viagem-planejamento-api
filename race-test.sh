#!/bin/sh
go test -race ./... -args config-location=./configs/config.json
