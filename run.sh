#!/bin/bash
OUTPUT="$(go build main.go 2>&1)"
if [[ ! -z ${OUTPUT} ]]; then
    echo ${OUTPUT}
else
    ./main printer.go
fi
