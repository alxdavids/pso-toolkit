#!/bin/bash

go test -timeout 900m -args -k 2048 -n 256
go test -timeout 900m -args -k 2048 -n 1024
go test -timeout 900m -args -k 2048 -n 4096
go test -timeout 900m -args -k 2048 -n 16384
go test -timeout 900m -args -k 2048 -n 65536
go test -timeout 900m -args -k 2048 -n 262144