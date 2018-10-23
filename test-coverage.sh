#!/bin/sh
go test ./... -coverprofile=/tmp/go-code-cover -args config-location=./configs/config.json
