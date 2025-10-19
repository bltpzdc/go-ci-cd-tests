#!/bin/bash
docker compose up -d
go test ./... -count=1
