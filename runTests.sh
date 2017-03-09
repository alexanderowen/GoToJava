#!/bin/bash
# Executes 'main.go' on each .go file in 'files', then
# attempts to compile output with javac. If javac reports
# errors, then they are echo'd 
# 
# Usage: './runTests.sh'

RED='\033[0;31m'
NC='\033[0m'
OUTPUT="$(go build main.go 2>&1)"
if [[ ! -z ${OUTPUT} ]]; then
    echo ${OUTPUT}
else
    #./main files/loop.go
    for src in `ls files/*.go`; do
        srcname=`basename $src`
        $(./main $src > t.java)
        LOCALOUTPUT=$(javac t.java 2>&1)
        if [[ ! -z ${LOCALOUTPUT} ]]; then
            echo -e "${RED}Failed $srcname${NC}"
            echo "$LOCALOUTPUT"
        fi
    done
fi
