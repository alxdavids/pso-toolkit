#!/bin/bash

$GOEXE test -timeout 1800m -args -k 2048 -m 16 -n 256
$GOEXE test -timeout 1800m -args -k 2048 -m 16 -n 1024
$GOEXE test -timeout 1800m -args -k 2048 -m 16 -n 4096
$GOEXE test -timeout 1800m -args -k 2048 -m 16 -n 16384
$GOEXE test -timeout 1800m -args -k 2048 -m 16 -n 65536
$GOEXE test -timeout 1800m -args -k 2048 -m 16 -n 262144