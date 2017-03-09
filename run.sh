#!/bin/bash
#
# Usage: "./run.sh path/to/file"

FILE=$1
OUTPUT="$(go build main.go 2>&1)"
if [[ ! -z ${OUTPUT} ]]; then
    echo ${OUTPUT}
else
    if [[ ! -z ${FILE} ]]; then
        ./main $FILE 
    else
        echo "No file specified"
    fi
fi
