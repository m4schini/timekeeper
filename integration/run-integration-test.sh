#!/bin/sh
cp ../database/schema.sql ./schema.sql

DOCKER_HOST=unix://$XDG_RUNTIME_DIR/podman/podman.sock go test ./... -json

