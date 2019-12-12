#!/bin/bash
cd cmd/switch-board-service && CGO_ENABLED=0 GOOS=linux go build -o switch-board-service -a -installsuffix cgo -ldflags '-w'