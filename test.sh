#!/bin/bash

$GOEXE test -timeout 900m -args -k 1024 -n 256
$GOEXE test -timeout 900m -args -k 1024 -n 1024
$GOEXE test -timeout 900m -args -k 1024 -n 4096
$GOEXE test -timeout 900m -args -k 1024 -n 16384
$GOEXE test -timeout 900m -args -k 1024 -n 65536
$GOEXE test -timeout 900m -args -k 1024 -n 262144