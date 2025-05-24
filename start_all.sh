#!/bin/bash

echo "Starting main application..."
go run cmd/main.go &

echo "Starting listener service..."
go run internal/script/listener/main.go &

echo "Both services started in the background."
wait 