#!/bin/bash

set -ex

cd slack-bot

# Dowloading packages
go mod tidy -compact=1.17

# running the slackbot
go run cmd/edith/main.go
