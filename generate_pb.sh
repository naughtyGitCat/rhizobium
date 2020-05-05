#!/usr/bin/env bash
protoc -I . rpc/rhizobium.proto --go_out=plugins=grpc:.
